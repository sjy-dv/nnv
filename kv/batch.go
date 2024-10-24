package kv

import (
	"fmt"
	"sync"

	"github.com/sjy-dv/nnv/pkg/diskhash"
	"github.com/sjy-dv/nnv/pkg/snowflake"
)

type Batch struct {
	db            *DB
	pendingWrites map[string]*LogRecord
	options       BatchOptions
	mu            sync.RWMutex
	committed     bool
	batchID       *snowflake.Node
}

// NewBatch creates a new Batch instance.
func (db *DB) NewBatch(options BatchOptions) *Batch {
	batch := &Batch{
		db:        db,
		options:   options,
		committed: false,
	}
	if !options.ReadOnly {
		batch.pendingWrites = make(map[string]*LogRecord)
		node, err := snowflake.NewNode(1)
		if err != nil {
			panic(fmt.Sprintf("snowflake.NewNode(1) failed: %v", err))
		}
		batch.batchID = node
	}
	batch.lock()
	return batch
}

func makeBatch() interface{} {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(fmt.Sprintf("snowflake.NewNode(1) failed: %v", err))
	}
	return &Batch{
		options: DefaultBatchOptions,
		batchID: node,
	}
}

func (b *Batch) init(rdonly, sync bool, disableWal bool, db *DB) *Batch {
	b.options.ReadOnly = rdonly
	b.options.Sync = sync
	b.options.DisableWal = disableWal
	b.db = db
	b.lock()
	return b
}

func (b *Batch) withPendingWrites() {
	b.pendingWrites = make(map[string]*LogRecord)
}

func (b *Batch) reset() {
	b.db = nil
	b.pendingWrites = nil
	b.committed = false
}

func (b *Batch) lock() {
	if b.options.ReadOnly {
		b.db.mu.RLock()
	} else {
		b.db.mu.Lock()
	}
}

func (b *Batch) unlock() {
	if b.options.ReadOnly {
		b.db.mu.RUnlock()
	} else {
		b.db.mu.Unlock()
	}
}

func (b *Batch) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}
	if b.db.closed {
		return ErrDBClosed
	}
	if b.options.ReadOnly {
		return ErrReadOnlyBatch
	}

	b.mu.Lock()

	b.pendingWrites[string(key)] = &LogRecord{
		Key:   key,
		Value: value,
		Type:  LogRecordNormal,
	}
	b.mu.Unlock()

	return nil
}

func (b *Batch) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, ErrKeyIsEmpty
	}
	if b.db.closed {
		return nil, ErrDBClosed
	}

	if b.pendingWrites != nil {
		b.mu.RLock()
		if record := b.pendingWrites[string(key)]; record != nil {
			if record.Type == LogRecordDeleted {
				b.mu.RUnlock()
				return nil, ErrKeyNotFound
			}
			b.mu.RUnlock()
			return record.Value, nil
		}
		b.mu.RUnlock()
	}

	tables := b.db.getMemTables()
	for _, table := range tables {
		deleted, value := table.get(key)
		if deleted {
			return nil, ErrKeyNotFound
		}
		if len(value) != 0 {
			return value, nil
		}
	}

	var value []byte
	var matchKey func(diskhash.Slot) (bool, error)
	if b.db.options.IndexType == Hash {
		matchKey = MatchKeyFunc(b.db, key, nil, &value)
	}

	position, err := b.db.index.Get(key, matchKey)
	if err != nil {
		return nil, err
	}

	if b.db.options.IndexType == Hash {
		if value == nil {
			return nil, ErrKeyNotFound
		}
		return value, nil
	}
	if position == nil {
		return nil, ErrKeyNotFound
	}
	record, err := b.db.vlog.read(position)
	if err != nil {
		return nil, err
	}
	return record.value, nil
}

func (b *Batch) Delete(key []byte) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}
	if b.db.closed {
		return ErrDBClosed
	}
	if b.options.ReadOnly {
		return ErrReadOnlyBatch
	}

	b.mu.Lock()
	b.pendingWrites[string(key)] = &LogRecord{
		Key:  key,
		Type: LogRecordDeleted,
	}
	b.mu.Unlock()

	return nil
}

func (b *Batch) Exist(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, ErrKeyIsEmpty
	}
	if b.db.closed {
		return false, ErrDBClosed
	}

	if b.pendingWrites != nil {
		b.mu.RLock()
		if record := b.pendingWrites[string(key)]; record != nil {
			b.mu.RUnlock()
			return record.Type != LogRecordDeleted, nil
		}
		b.mu.RUnlock()
	}

	tables := b.db.getMemTables()
	for _, table := range tables {
		deleted, value := table.get(key)
		if deleted {
			return false, nil
		}
		if len(value) != 0 {
			return true, nil
		}
	}

	var value []byte
	var matchKeyFunc func(diskhash.Slot) (bool, error)
	if b.db.options.IndexType == Hash {
		matchKeyFunc = MatchKeyFunc(b.db, key, nil, &value)
	}
	pos, err := b.db.index.Get(key, matchKeyFunc)
	if err != nil {
		return false, err
	}
	if b.db.options.IndexType == Hash {
		return value != nil, nil
	}
	return pos != nil, nil
}

func (b *Batch) Commit() error {
	defer b.unlock()
	if b.db.closed {
		return ErrDBClosed
	}

	if b.options.ReadOnly || len(b.pendingWrites) == 0 {
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.committed {
		return ErrBatchCommitted
	}

	if err := b.db.waitMemtableSpace(); err != nil {
		return err
	}
	batchID := b.batchID.Generate()

	err := b.db.activeMem.putBatch(b.pendingWrites, batchID, b.options.WriteOptions)
	if err != nil {
		return err
	}

	b.committed = true
	return nil
}
