package wal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

const (
	initialSegmentFileID = 1
)

var (
	ErrValueTooLarge       = errors.New("the data size can't larger than segment size")
	ErrPendingSizeTooLarge = errors.New("the upper bound of pendingWrites can't larger than segment size")
)

type WAL struct {
	activeSegment     *segment               // active segment file, used for new incoming writes.
	olderSegments     map[SegmentID]*segment // older segment files, only used for read.
	options           Options
	mu                sync.RWMutex
	bytesWrite        uint32
	renameIds         []SegmentID
	pendingWrites     [][]byte
	pendingSize       int64
	pendingWritesLock sync.Mutex
}

type Reader struct {
	segmentReaders []*segmentReader
	currentReader  int
}

func Open(options Options) (*WAL, error) {
	if !strings.HasPrefix(options.SegmentFileExt, ".") {
		return nil, fmt.Errorf("segment file extension must start with '.'")
	}
	wal := &WAL{
		options:       options,
		olderSegments: make(map[SegmentID]*segment),
		pendingWrites: make([][]byte, 0),
	}

	// create the directory if not exists.
	if err := os.MkdirAll(options.DirPath, os.ModePerm); err != nil {
		return nil, err
	}

	// iterate the dir and open all segment files.
	entries, err := os.ReadDir(options.DirPath)
	if err != nil {
		return nil, err
	}

	// get all segment file ids.
	var segmentIDs []int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		var id int
		_, err := fmt.Sscanf(entry.Name(), "%d"+options.SegmentFileExt, &id)
		if err != nil {
			continue
		}
		segmentIDs = append(segmentIDs, id)
	}

	// empty directory, just initialize a new segment file.
	if len(segmentIDs) == 0 {
		segment, err := openSegmentFile(options.DirPath, options.SegmentFileExt,
			initialSegmentFileID)
		if err != nil {
			return nil, err
		}
		wal.activeSegment = segment
	} else {
		// open the segment files in order, get the max one as the active segment file.
		sort.Ints(segmentIDs)

		for i, segId := range segmentIDs {
			segment, err := openSegmentFile(options.DirPath, options.SegmentFileExt,
				uint32(segId))
			if err != nil {
				return nil, err
			}
			if i == len(segmentIDs)-1 {
				wal.activeSegment = segment
			} else {
				wal.olderSegments[segment.id] = segment
			}
		}
	}

	return wal, nil
}

func SegmentFileName(dirPath string, extName string, id SegmentID) string {
	return filepath.Join(dirPath, fmt.Sprintf("%09d"+extName, id))
}

func (wal *WAL) OpenNewActiveSegment() error {
	wal.mu.Lock()
	defer wal.mu.Unlock()
	// sync the active segment file.
	if err := wal.activeSegment.Sync(); err != nil {
		return err
	}
	// create a new segment file and set it as the active one.
	segment, err := openSegmentFile(wal.options.DirPath, wal.options.SegmentFileExt,
		wal.activeSegment.id+1)
	if err != nil {
		return err
	}
	wal.olderSegments[wal.activeSegment.id] = wal.activeSegment
	wal.activeSegment = segment
	return nil
}

// ActiveSegmentID returns the id of the active segment file.
func (wal *WAL) ActiveSegmentID() SegmentID {
	wal.mu.RLock()
	defer wal.mu.RUnlock()

	return wal.activeSegment.id
}

func (wal *WAL) IsEmpty() bool {
	wal.mu.RLock()
	defer wal.mu.RUnlock()

	return len(wal.olderSegments) == 0 && wal.activeSegment.Size() == 0
}

func (wal *WAL) SetIsStartupTraversal(v bool) {
	for _, seg := range wal.olderSegments {
		seg.isStartupTraversal = v
	}
	wal.activeSegment.isStartupTraversal = v
}

func (wal *WAL) NewReaderWithMax(segId SegmentID) *Reader {
	wal.mu.RLock()
	defer wal.mu.RUnlock()

	// get all segment readers.
	var segmentReaders []*segmentReader
	for _, segment := range wal.olderSegments {
		if segId == 0 || segment.id <= segId {
			reader := segment.NewReader()
			segmentReaders = append(segmentReaders, reader)
		}
	}
	if segId == 0 || wal.activeSegment.id <= segId {
		reader := wal.activeSegment.NewReader()
		segmentReaders = append(segmentReaders, reader)
	}

	// sort the segment readers by segment id.
	sort.Slice(segmentReaders, func(i, j int) bool {
		return segmentReaders[i].segment.id < segmentReaders[j].segment.id
	})

	return &Reader{
		segmentReaders: segmentReaders,
		currentReader:  0,
	}
}

func (wal *WAL) NewReaderWithStart(startPos *ChunkPosition) (*Reader, error) {
	if startPos == nil {
		return nil, errors.New("start position is nil")
	}
	wal.mu.RLock()
	defer wal.mu.RUnlock()

	reader := wal.NewReader()
	for {
		// skip the segment readers whose id is less than the given position's segment id.
		if reader.CurrentSegmentId() < startPos.SegmentId {
			reader.SkipCurrentSegment()
			continue
		}
		// skip the chunk whose position is less than the given position.
		currentPos := reader.CurrentChunkPosition()
		if currentPos.BlockNumber >= startPos.BlockNumber &&
			currentPos.ChunkOffset >= startPos.ChunkOffset {
			break
		}
		// call Next to find again.
		if _, _, err := reader.Next(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return reader, nil
}

func (wal *WAL) NewReader() *Reader {
	return wal.NewReaderWithMax(0)
}

func (r *Reader) Next() ([]byte, *ChunkPosition, error) {
	if r.currentReader >= len(r.segmentReaders) {
		return nil, nil, io.EOF
	}

	data, position, err := r.segmentReaders[r.currentReader].Next()
	if err == io.EOF {
		r.currentReader++
		return r.Next()
	}
	return data, position, err
}

func (r *Reader) SkipCurrentSegment() {
	r.currentReader++
}

func (r *Reader) CurrentSegmentId() SegmentID {
	return r.segmentReaders[r.currentReader].segment.id
}

func (r *Reader) CurrentChunkPosition() *ChunkPosition {
	reader := r.segmentReaders[r.currentReader]
	return &ChunkPosition{
		SegmentId:   reader.segment.id,
		BlockNumber: reader.blockNumber,
		ChunkOffset: reader.chunkOffset,
	}
}

func (wal *WAL) ClearPendingWrites() {
	wal.pendingWritesLock.Lock()
	defer wal.pendingWritesLock.Unlock()

	wal.pendingSize = 0
	wal.pendingWrites = wal.pendingWrites[:0]
}

func (wal *WAL) PendingWrites(data []byte) {
	wal.pendingWritesLock.Lock()
	defer wal.pendingWritesLock.Unlock()

	size := wal.maxDataWriteSize(int64(len(data)))
	wal.pendingSize += size
	wal.pendingWrites = append(wal.pendingWrites, data)
}

func (wal *WAL) rotateActiveSegment() error {
	if err := wal.activeSegment.Sync(); err != nil {
		return err
	}
	wal.bytesWrite = 0
	segment, err := openSegmentFile(wal.options.DirPath, wal.options.SegmentFileExt,
		wal.activeSegment.id+1)
	if err != nil {
		return err
	}
	wal.olderSegments[wal.activeSegment.id] = wal.activeSegment
	wal.activeSegment = segment
	return nil
}

func (wal *WAL) WriteAll() ([]*ChunkPosition, error) {
	if len(wal.pendingWrites) == 0 {
		return make([]*ChunkPosition, 0), nil
	}

	wal.mu.Lock()
	defer func() {
		wal.ClearPendingWrites()
		wal.mu.Unlock()
	}()

	// if the pending size is still larger than segment size, return error
	if wal.pendingSize > wal.options.SegmentSize {
		return nil, ErrPendingSizeTooLarge
	}

	// if the active segment file is full, sync it and create a new one.
	if wal.activeSegment.Size()+wal.pendingSize > wal.options.SegmentSize {
		if err := wal.rotateActiveSegment(); err != nil {
			return nil, err
		}
	}

	// write all data to the active segment file.
	positions, err := wal.activeSegment.writeAll(wal.pendingWrites)
	if err != nil {
		return nil, err
	}

	return positions, nil
}

func (wal *WAL) Write(data []byte) (*ChunkPosition, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()
	if int64(len(data))+chunkHeaderSize > wal.options.SegmentSize {
		return nil, ErrValueTooLarge
	}
	// if the active segment file is full, sync it and create a new one.
	if wal.isFull(int64(len(data))) {
		if err := wal.rotateActiveSegment(); err != nil {
			return nil, err
		}
	}

	// write the data to the active segment file.
	position, err := wal.activeSegment.Write(data)
	if err != nil {
		return nil, err
	}

	// update the bytesWrite field.
	wal.bytesWrite += position.ChunkSize

	// sync the active segment file if needed.
	var needSync = wal.options.Sync
	if !needSync && wal.options.BytesPerSync > 0 {
		needSync = wal.bytesWrite >= wal.options.BytesPerSync
	}
	if needSync {
		if err := wal.activeSegment.Sync(); err != nil {
			return nil, err
		}
		wal.bytesWrite = 0
	}

	return position, nil
}

// Read reads the data from the WAL according to the given position.
func (wal *WAL) Read(pos *ChunkPosition) ([]byte, error) {
	wal.mu.RLock()
	defer wal.mu.RUnlock()

	// find the segment file according to the position.
	var segment *segment
	if pos.SegmentId == wal.activeSegment.id {
		segment = wal.activeSegment
	} else {
		segment = wal.olderSegments[pos.SegmentId]
	}

	if segment == nil {
		return nil, fmt.Errorf("segment file %d%s not found", pos.SegmentId, wal.options.SegmentFileExt)
	}

	// read the data from the segment file.
	return segment.Read(pos.BlockNumber, pos.ChunkOffset)
}

// Close closes the WAL.
func (wal *WAL) Close() error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// close all segment files.
	for _, segment := range wal.olderSegments {
		if err := segment.Close(); err != nil {
			return err
		}
		wal.renameIds = append(wal.renameIds, segment.id)
	}
	wal.olderSegments = nil

	wal.renameIds = append(wal.renameIds, wal.activeSegment.id)
	// close the active segment file.
	return wal.activeSegment.Close()
}

// Delete deletes all segment files of the WAL.
func (wal *WAL) Delete() error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// delete all segment files.
	for _, segment := range wal.olderSegments {
		if err := segment.Remove(); err != nil {
			return err
		}
	}
	wal.olderSegments = nil

	// delete the active segment file.
	return wal.activeSegment.Remove()
}

// Sync syncs the active segment file to stable storage like disk.
func (wal *WAL) Sync() error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	return wal.activeSegment.Sync()
}

func (wal *WAL) RenameFileExt(ext string) error {
	if !strings.HasPrefix(ext, ".") {
		return fmt.Errorf("segment file extension must start with '.'")
	}
	wal.mu.Lock()
	defer wal.mu.Unlock()

	renameFile := func(id SegmentID) error {
		oldName := SegmentFileName(wal.options.DirPath, wal.options.SegmentFileExt, id)
		newName := SegmentFileName(wal.options.DirPath, ext, id)
		return os.Rename(oldName, newName)
	}

	for _, id := range wal.renameIds {
		if err := renameFile(id); err != nil {
			return err
		}
	}

	wal.options.SegmentFileExt = ext
	return nil
}

func (wal *WAL) isFull(delta int64) bool {
	return wal.activeSegment.Size()+wal.maxDataWriteSize(delta) > wal.options.SegmentSize
}

func (wal *WAL) maxDataWriteSize(size int64) int64 {
	return chunkHeaderSize + size + (size/blockSize+1)*chunkHeaderSize
}