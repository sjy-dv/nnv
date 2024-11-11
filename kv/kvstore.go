// Licensed to sjy-dv under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. sjy-dv licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package kv

import (
	"context"
	"errors"
	"fmt"
	"io"

	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v4/y"
	"github.com/gofrs/flock"
	"github.com/rs/zerolog/log"
	"github.com/sjy-dv/nnv/kv/wal"
	"github.com/sjy-dv/nnv/pkg/diskhash"
	"golang.org/x/sync/errgroup"
)

const (
	fileLockName = "FLOCK"
)

type DB struct {
	activeMem *memtable      // Active memtable for writing.
	immuMems  []*memtable    // Immutable memtables, waiting to be flushed to disk.
	index     Index          // index is multi-partition indexes to store key and chunk position.
	vlog      *valueLog      // vlog is the value log.
	fileLock  *flock.Flock   // fileLock to prevent multiple processes from using the same database directory.
	flushChan chan *memtable // flushChan is used to notify the flush goroutine to flush memtable to disk.
	flushLock sync.Mutex     // flushLock is to prevent flush running while compaction doesn't occur
	mu        sync.RWMutex
	closed    bool
	closeChan chan struct{}
	options   Options
	batchPool sync.Pool // batchPool is a pool of batch, to reduce the cost of memory allocation.
}

func Open(options Options) (*DB, error) {
	// check whether all options are valid
	if err := validateOptions(&options); err != nil {
		return nil, err
	}

	// create data directory if not exist
	if _, err := os.Stat(options.DirPath); err != nil {
		if err = os.MkdirAll(options.DirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// create file lock, prevent multiple processes from using the same database directory
	fileLock := flock.New(filepath.Join(options.DirPath, fileLockName))
	hold, err := fileLock.TryLock()
	if err != nil {
		return nil, err
	}
	if !hold {
		return nil, ErrDatabaseIsUsing
	}

	// open all memtables
	memtables, err := openAllMemtables(options)
	if err != nil {
		return nil, err
	}

	// open index
	index, err := openIndex(indexOptions{
		indexType:       options.IndexType,
		dirPath:         options.DirPath,
		partitionNum:    options.PartitionNum,
		keyHashFunction: options.KeyHashFunction,
	})
	if err != nil {
		return nil, err
	}

	// open value log
	vlog, err := openValueLog(valueLogOptions{
		dirPath:           options.DirPath,
		segmentSize:       options.ValueLogFileSize,
		blockCache:        options.BlockCache,
		partitionNum:      uint32(options.PartitionNum),
		hashKeyFunction:   options.KeyHashFunction,
		compactBatchCount: options.CompactBatchCount,
	})
	if err != nil {
		return nil, err
	}

	db := &DB{
		activeMem: memtables[len(memtables)-1],
		immuMems:  memtables[:len(memtables)-1],
		index:     index,
		vlog:      vlog,
		fileLock:  fileLock,
		flushChan: make(chan *memtable, options.MemtableNums-1),
		closeChan: make(chan struct{}),
		options:   options,
		batchPool: sync.Pool{New: makeBatch},
	}

	// if there are some immutable memtables when opening the database, flush them to disk
	if len(db.immuMems) > 0 {
		for _, table := range db.immuMems {
			db.flushMemtable(table)
		}
	}

	// start flush memtables goroutine asynchronously,
	// memtables with new coming writes will be flushed to disk if the active memtable is full.
	go db.listenMemtableFlush()

	return db, nil
}

// Close the database, close all data files and release file lock.
// Set the closed flag to true.
// The DB instance cannot be used after closing.
func (db *DB) Close() error {
	log.Debug().Msg("Attempting to acquire lock for closing DB")
	db.mu.Lock()
	defer db.mu.Unlock()
	log.Debug().Msg("Lock acquired, starting close process for DB")

	log.Debug().Msg("Closing flush channel")

	close(db.flushChan)

	closeTimeout := time.After(1 * time.Second)

	log.Debug().Msg("Waiting for close channel to signal")
	select {
	case <-db.closeChan:
		log.Debug().Msg("Close channel signaled successfully")
	case <-closeTimeout:
		log.Warn().Msg("Timeout reached while waiting for close channel signal")
	}
	log.Debug().Msg("Close channel signaled")

	// Close all immutable memtables
	for i, table := range db.immuMems {
		log.Debug().Msgf("Closing immutable memtable #%d", i)
		if err := table.close(); err != nil {
			log.Warn().Err(err).Msgf("Failed to close immutable memtable #%d", i)
			return err
		}
	}
	log.Debug().Msg("All immutable memtables closed successfully")

	// Close active memtable
	log.Debug().Msg("Closing active memtable")
	if err := db.activeMem.close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close active memtable")
		return err
	}
	log.Debug().Msg("Active memtable closed successfully")

	// Close index
	log.Debug().Msg("Closing index")
	if err := db.index.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close index")
		return err
	}
	log.Debug().Msg("Index closed successfully")

	// Close value log
	log.Debug().Msg("Closing value log")
	if err := db.vlog.close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close value log")
		return err
	}
	log.Debug().Msg("Value log closed successfully")

	// Release file lock
	log.Debug().Msg("Releasing file lock")
	if err := db.fileLock.Unlock(); err != nil {
		log.Warn().Err(err).Msg("Failed to release file lock")
		return err
	}
	log.Debug().Msg("File lock released successfully")

	db.closed = true
	log.Debug().Msg("DB closed successfully")
	return nil
}

// Sync all data files to the underlying storage.
func (db *DB) Sync() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// sync all wal of memtables
	for _, table := range db.immuMems {
		if err := table.sync(); err != nil {
			return err
		}
	}
	if err := db.activeMem.sync(); err != nil {
		return err
	}
	// sync index
	if err := db.index.Sync(); err != nil {
		return err
	}
	// sync value log
	if err := db.vlog.sync(); err != nil {
		return err
	}

	return nil
}

// Put put with defaultWriteOptions.
func (db *DB) Put(key []byte, value []byte) error {
	return db.PutWithOptions(key, value, DefaultWriteOptions)
}

// PutWithOptions a key-value pair into the database.
// Actually, it will open a new batch and commit it.
// You can think the batch has only one Put operation.
func (db *DB) PutWithOptions(key []byte, value []byte, options WriteOptions) error {
	batch, ok := db.batchPool.Get().(*Batch)
	if !ok {
		panic("batchPoll.Get failed")
	}
	batch.options.WriteOptions = options
	defer func() {
		batch.reset()
		db.batchPool.Put(batch)
	}()
	// This is a single put operation, we can set Sync to false.
	// Because the data will be written to the WAL,
	// and the WAL file will be synced to disk according to the DB options.
	batch.init(false, false, false, db).withPendingWrites()
	if err := batch.Put(key, value); err != nil {
		batch.unlock()
		return err
	}
	return batch.Commit()
}

// Get the value of the specified key from the database.
// Actually, it will open a new batch and commit it.
// You can think the batch has only one Get operation.
func (db *DB) Get(key []byte) ([]byte, error) {
	batch, ok := db.batchPool.Get().(*Batch)
	if !ok {
		panic("batchPoll.Get failed")
	}
	batch.init(true, false, true, db)
	defer func() {
		_ = batch.Commit()
		batch.reset()
		db.batchPool.Put(batch)
	}()
	return batch.Get(key)
}

// Delete delete with defaultWriteOptions.
func (db *DB) Delete(key []byte) error {
	return db.DeleteWithOptions(key, DefaultWriteOptions)
}

// DeleteWithOptions the specified key from the database.
// Actually, it will open a new batch and commit it.
// You can think the batch has only one Delete operation.
func (db *DB) DeleteWithOptions(key []byte, options WriteOptions) error {
	batch, ok := db.batchPool.Get().(*Batch)
	if !ok {
		panic("batchPoll.Get failed")
	}
	batch.options.WriteOptions = options
	defer func() {
		batch.reset()
		db.batchPool.Put(batch)
	}()
	// This is a single delete operation, we can set Sync to false.
	// Because the data will be written to the WAL,
	// and the WAL file will be synced to disk according to the DB options.
	batch.init(false, false, false, db).withPendingWrites()
	if err := batch.Delete(key); err != nil {
		batch.unlock()
		return err
	}
	return batch.Commit()
}

// Exist checks if the specified key exists in the database.
// Actually, it will open a new batch and commit it.
// You can think the batch has only one Exist operation.
func (db *DB) Exist(key []byte) (bool, error) {
	batch, ok := db.batchPool.Get().(*Batch)
	if !ok {
		panic("batchPoll.Get failed")
	}
	batch.init(true, false, true, db)
	defer func() {
		_ = batch.Commit()
		batch.reset()
		db.batchPool.Put(batch)
	}()
	return batch.Exist(key)
}

// validateOptions validates the given options.
func validateOptions(options *Options) error {
	if options.DirPath == "" {
		return ErrDBDirectoryISEmpty
	}
	if options.MemtableSize <= 0 {
		options.MemtableSize = DefaultOptions.MemtableSize
	}
	if options.MemtableNums <= 0 {
		options.MemtableNums = DefaultOptions.MemtableNums
	}
	if options.PartitionNum <= 0 {
		options.PartitionNum = DefaultOptions.PartitionNum
	}
	if options.ValueLogFileSize <= 0 {
		options.ValueLogFileSize = DefaultOptions.ValueLogFileSize
	}
	// assure ValueLogFileSize >= MemtableSize
	if options.ValueLogFileSize < int64(options.MemtableSize) {
		options.ValueLogFileSize = int64(options.MemtableSize)
	}
	return nil
}

// get all memtables, including active memtable and immutable memtables.
// must be called with db.mu held.
func (db *DB) getMemTables() []*memtable {
	var tables []*memtable
	tables = append(tables, db.activeMem)

	last := len(db.immuMems) - 1
	for i := range db.immuMems {
		tables = append(tables, db.immuMems[last-i])
	}

	return tables
}

// waitMemtableSpace waits for space in the memtable.
// If the active memtable is full, it will be flushed to disk by the background goroutine.
// But if the flush speed is slower than the write speed, there may be no space in the memtable.
// So the write operation will wait for space in the memtable, and the timeout is specified by WaitMemSpaceTimeout.
func (db *DB) waitMemtableSpace() error {
	if !db.activeMem.isFull() {
		return nil
	}

	timer := time.NewTimer(db.options.WaitMemSpaceTimeout)
	defer timer.Stop()
	select {
	case db.flushChan <- db.activeMem:
		db.immuMems = append(db.immuMems, db.activeMem)
		options := db.activeMem.options
		options.tableID++
		// open a new memtable for writing
		table, err := openMemtable(options)
		if err != nil {
			return err
		}
		db.activeMem = table
	case <-timer.C:
		return ErrWaitMemtableSpaceTimeOut
	}

	return nil
}

// flushMemtable flushes the specified memtable to disk.
// Following steps will be done:
// 1. Iterate all records in memtable, divide them into deleted keys and log records.
// 2. Write the log records to value log, get the positions of keys.
// 3. Write all keys and positions to index.
// 4. Delete the deleted keys from index.
// 5. Delete the wal.
func (db *DB) flushMemtable(table *memtable) {
	db.flushLock.Lock()
	defer db.flushLock.Unlock()
	sklIter := table.skl.NewIterator()
	var deletedKeys [][]byte
	var logRecords []*ValueLogRecord

	// iterate all records in memtable, divide them into deleted keys and log records
	for sklIter.SeekToFirst(); sklIter.Valid(); sklIter.Next() {
		key, valueStruct := y.ParseKey(sklIter.Key()), sklIter.Value()
		if valueStruct.Meta == LogRecordDeleted {
			deletedKeys = append(deletedKeys, key)
		} else {
			logRecord := ValueLogRecord{key: key, value: valueStruct.Value}
			logRecords = append(logRecords, &logRecord)
		}
	}
	_ = sklIter.Close()

	// write to value log, get the positions of keys
	keyPos, err := db.vlog.writeBatch(logRecords)
	if err != nil {
		log.Warn().Msgf("vlog writeBatch failed:%v", err)
		return
	}

	// sync the value log
	if err = db.vlog.sync(); err != nil {
		log.Warn().Msgf("vlog sync failed:%v", err)
		return
	}

	// write all keys and positions to index
	var putMatchKeys []diskhash.MatchKeyFunc
	if db.options.IndexType == Hash && len(keyPos) > 0 {
		putMatchKeys = make([]diskhash.MatchKeyFunc, len(keyPos))
		for i := range putMatchKeys {
			putMatchKeys[i] = MatchKeyFunc(db, keyPos[i].key, nil, nil)
		}
	}
	if err = db.index.PutBatch(keyPos, putMatchKeys...); err != nil {
		log.Warn().Msgf("index PutBatch failed:", err)
		return
	}
	// delete the deleted keys from index
	var deleteMatchKeys []diskhash.MatchKeyFunc
	if db.options.IndexType == Hash && len(deletedKeys) > 0 {
		deleteMatchKeys = make([]diskhash.MatchKeyFunc, len(deletedKeys))
		for i := range deleteMatchKeys {
			deleteMatchKeys[i] = MatchKeyFunc(db, deletedKeys[i], nil, nil)
		}
	}
	if err = db.index.DeleteBatch(deletedKeys, deleteMatchKeys...); err != nil {
		log.Warn().Msgf("index DeleteBatch failed:%v", err)
		return
	}
	// sync the index
	if err = db.index.Sync(); err != nil {
		log.Warn().Msgf("index sync failed:%v", err)
		return
	}

	// delete the wal
	if err = table.deleteWAl(); err != nil {
		log.Warn().Msgf("delete wal failed:%v", err)
		return
	}

	// delete old memtable kept in memory
	db.mu.Lock()
	if table == db.activeMem {
		options := db.activeMem.options
		options.tableID++
		// open a new memtable for writing
		table, err = openMemtable(options)
		if err != nil {
			panic("flush activate memtable wrong")
		}
		db.activeMem = table
	} else {
		if len(db.immuMems) == 1 {
			db.immuMems = db.immuMems[:0]
		} else {
			db.immuMems = db.immuMems[1:]
		}
	}

	db.mu.Unlock()
}

func (db *DB) listenMemtableFlush() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		select {
		case table, ok := <-db.flushChan:
			if ok {
				db.flushMemtable(table)
			} else {
				db.closeChan <- struct{}{}
				return
			}
		case <-sig:
			return
		}
	}
}

// Compact will iterate all values in vlog, and write the valid values to a new vlog file.
// Then replace the old vlog file with the new one, and delete the old one.
//
//nolint:gocognit
func (db *DB) Compact() error {
	db.flushLock.Lock()
	defer db.flushLock.Unlock()

	openVlogFile := func(part int, ext string) *wal.WAL {
		walFile, err := wal.Open(wal.Options{
			DirPath:        db.vlog.options.dirPath,
			SegmentSize:    db.vlog.options.segmentSize,
			SegmentFileExt: fmt.Sprintf(ext, part),
			BlockCache:     db.vlog.options.blockCache,
			Sync:           false, // we will sync manually
			BytesPerSync:   0,     // the same as Sync
		})
		if err != nil {
			_ = walFile.Delete()
			panic(err)
		}
		return walFile
	}

	g, _ := errgroup.WithContext(context.Background())
	for i := 0; i < int(db.vlog.options.partitionNum); i++ {
		part := i
		g.Go(func() error {
			newVlogFile := openVlogFile(part, tempValueLogFileExt)

			validRecords := make([]*ValueLogRecord, 0, db.vlog.options.compactBatchCount)
			reader := db.vlog.walFiles[part].NewReader()
			count := 0
			// iterate all records in wal, find the valid records
			for {
				count++
				chunk, pos, err := reader.Next()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					_ = newVlogFile.Delete()
					return err
				}

				record := decodeValueLogRecord(chunk)
				var hashTableKeyPos *KeyPosition
				var matchKey func(diskhash.Slot) (bool, error)
				if db.options.IndexType == Hash {
					matchKey = MatchKeyFunc(db, record.key, &hashTableKeyPos, nil)
				}
				keyPos, err := db.index.Get(record.key, matchKey)
				if err != nil {
					_ = newVlogFile.Delete()
					return err
				}

				if db.options.IndexType == Hash {
					keyPos = hashTableKeyPos
				}

				if keyPos == nil {
					continue
				}
				if keyPos.partition == uint32(part) && reflect.DeepEqual(keyPos.position, pos) {
					validRecords = append(validRecords, record)
				}

				if count%db.vlog.options.compactBatchCount == 0 {
					err = db.rewriteValidRecords(newVlogFile, validRecords, part)
					if err != nil {
						_ = newVlogFile.Delete()
						return err
					}
					validRecords = validRecords[:0]
				}
			}

			if len(validRecords) > 0 {
				err := db.rewriteValidRecords(newVlogFile, validRecords, part)
				if err != nil {
					_ = newVlogFile.Delete()
					return err
				}
			}

			// replace the wal with the new one.
			_ = db.vlog.walFiles[part].Delete()
			_ = newVlogFile.Close()
			if err := newVlogFile.RenameFileExt(fmt.Sprintf(valueLogFileExt, part)); err != nil {
				return err
			}
			db.vlog.walFiles[part] = openVlogFile(part, valueLogFileExt)

			return nil
		})
	}

	return g.Wait()
}

func (db *DB) rewriteValidRecords(walFile *wal.WAL, validRecords []*ValueLogRecord, part int) error {
	for _, record := range validRecords {
		walFile.PendingWrites(encodeValueLogRecord(record))
	}

	walChunkPositions, err := walFile.WriteAll()
	if err != nil {
		return err
	}

	positions := make([]*KeyPosition, len(walChunkPositions))
	for i, walChunkPosition := range walChunkPositions {
		positions[i] = &KeyPosition{
			key:       validRecords[i].key,
			partition: uint32(part),
			position:  walChunkPosition,
		}
	}
	matchKeys := make([]diskhash.MatchKeyFunc, len(positions))
	if db.options.IndexType == Hash {
		for i := range matchKeys {
			matchKeys[i] = MatchKeyFunc(db, positions[i].key, nil, nil)
		}
	}
	return db.index.PutBatch(positions, matchKeys...)
}
