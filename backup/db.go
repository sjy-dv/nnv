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

package backup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/gofrs/uuid/v5"
	"github.com/sjy-dv/nnv/backup/document"
	"github.com/sjy-dv/nnv/backup/index"
	"github.com/sjy-dv/nnv/backup/internal"
	"github.com/sjy-dv/nnv/backup/query"
	"github.com/sjy-dv/nnv/backup/store"
	"github.com/sjy-dv/nnv/backup/store/bbolt"
)

// Collection creation errors
var (
	ErrCollectionExist    = errors.New("collection already exist")
	ErrCollectionNotExist = errors.New("no such collection")

	ErrIndexExist    = errors.New("index already exist")
	ErrIndexNotExist = errors.New("no such index")

	ErrDocumentNotExist = errors.New("no such document")
	ErrDuplicateKey     = errors.New("duplicate key")
)

type docConsumer func(doc *document.Document) error

// DB represents the entry point of each nnlog database.
type DB struct {
	store  store.Store
	closed uint32
}

type collectionMetadata struct {
	Size    int
	Indexes []index.Info
}

// CreateCollection creates a new empty collection with the given name.
func (db *DB) CreateCollection(name string) error {
	tx, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ok, err := db.hasCollection(name, tx)
	if err != nil {
		return err
	}

	if ok {
		return ErrCollectionExist
	}

	meta := &collectionMetadata{Size: 0}
	if err := db.saveCollectionMetadata(name, meta, tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) CreateCollectionByQuery(name string, q *query.Query) error {
	err := db.CreateCollection(name)
	if err != nil {
		return err
	}
	docs, err := db.FindAll(q)
	if err != nil {
		return err
	}
	if len(docs) == 0 { // just an empty collection
		return nil
	} else {
		return db.Insert(name, docs...)
	}
}

func (db *DB) saveCollectionMetadata(collection string, meta *collectionMetadata, tx store.Tx) error {
	rawMeta, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return tx.Set([]byte(getCollectionKey(collection)), rawMeta)
}

func (db *DB) hasCollection(name string, tx store.Tx) (bool, error) {
	value, err := tx.Get([]byte(getCollectionKey(name)))
	return value != nil, err
}

func getCollectionKey(name string) string {
	return getCollectionKeyPrefix() + name
}

func getCollectionKeyPrefix() string {
	return "coll:"
}

// DropCollection removes the collection with the given name, deleting any content on disk.
func (db *DB) DropCollection(name string) error {
	tx, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := db.deleteAll(tx, name); err != nil {
		return err
	}

	if err := tx.Delete([]byte(getCollectionKey(name))); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) deleteAll(tx store.Tx, collName string) error {
	return db.replaceDocs(tx, query.NewQuery(collName), func(_ *document.Document) *document.Document {
		return nil
	})
}

// HasCollection returns true if and only if the database contains a collection with the given name.
func (db *DB) HasCollection(name string) (bool, error) {
	txn, err := db.store.Begin(false)
	if err != nil {
		return false, err
	}
	defer txn.Rollback()
	return db.hasCollection(name, txn)
}

func NewObjectId() string {
	objId, _ := uuid.NewV4()
	return objId.String()
}

// Insert adds the supplied documents to a collection.
func (db *DB) Insert(collectionName string, docs ...*document.Document) error {
	for _, doc := range docs {
		if !doc.Has(document.ObjectIdField) || doc.Get(document.ObjectIdField) == "" {
			objectId := NewObjectId()
			doc.Set(document.ObjectIdField, objectId)
		}
	}

	tx, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	meta, err := db.getCollectionMeta(collectionName, tx)
	if err != nil {
		return err
	}

	indexes := db.getIndexes(tx, collectionName, meta)

	for _, doc := range docs {
		if err := db.addDocToIndexes(tx, indexes, doc); err != nil {
			return err
		}

		key := []byte(getDocumentKey(collectionName, doc.ObjectId()))
		value, err := tx.Get(key)
		if err != nil {
			return err
		}

		if value != nil {
			return ErrDuplicateKey
		}

		if err := saveDocument(doc, key, tx); err != nil {
			return err
		}
	}

	meta.Size += len(docs)
	if err := db.saveCollectionMetadata(collectionName, meta, tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) getIndexes(tx store.Tx, collection string, meta *collectionMetadata) []index.Index {
	indexes := make([]index.Index, 0)

	for _, info := range meta.Indexes {
		indexes = append(indexes, index.CreateIndex(collection, info.Field, info.Type, tx))
	}
	return indexes
}

func saveDocument(doc *document.Document, key []byte, tx store.Tx) error {
	if err := document.Validate(doc); err != nil {
		return err
	}

	data, err := document.Encode(doc)
	if err != nil {
		return err
	}
	return tx.Set(key, data)
}

func (db *DB) addDocToIndexes(tx store.Tx, indexes []index.Index, doc *document.Document) error {
	// update indexes
	for _, idx := range indexes {
		fieldVal := doc.Get(idx.Field()) // missing fields are treated as null

		err := idx.Add(doc.ObjectId(), fieldVal, doc.TTL())
		if err != nil {
			return err
		}
	}
	return nil
}

func getDocumentKey(collection string, id string) string {
	return getDocumentKeyPrefix(collection) + id
}

func getDocumentKeyPrefix(collection string) string {
	return "c:" + collection + ";" + "d:"
}

func (db *DB) getCollectionMeta(collection string, tx store.Tx) (*collectionMetadata, error) {
	value, err := tx.Get([]byte(getCollectionKey(collection)))
	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, ErrCollectionNotExist
	}

	m := &collectionMetadata{}
	err = json.Unmarshal(value, m)
	return m, err
}

// Save or update a document, If you pass in a custom struct instead of a Document object,
// it is recommended to specify the _id field using struct tags.
func (db *DB) Save(collectionName string, data interface{}) error {
	doc := document.NewDocumentOf(data)
	if !doc.Has(document.ObjectIdField) || doc.Get(document.ObjectIdField) == "" {
		return db.Insert(collectionName, doc)
	}
	return db.ReplaceById(collectionName, doc.ObjectId(), doc)
}

// InsertOne inserts a single document to an existing collection. It returns the id of the inserted document.
func (db *DB) InsertOne(collectionName string, doc *document.Document) (string, error) {
	err := db.Insert(collectionName, doc)
	return doc.ObjectId(), err
}

// Open opens a new nnlog database on the supplied path. If such a folder doesn't exist, it is automatically createdocument.
func Open(dir string) (*DB, error) {
	dataStore, err := bbolt.Open(dir)
	if err != nil {
		return nil, err
	}
	return OpenWithStore(dataStore)
}

// OpenWithStore opens a new nnlog database using the provided store.
func OpenWithStore(store store.Store) (*DB, error) {
	return &DB{store: store}, nil
}

// Close releases all the resources and closes the database. After the call, the instance will no more be usable.
func (db *DB) Close() error {
	if atomic.CompareAndSwapUint32(&db.closed, 0, 1) {
		return db.store.Close()
	}
	return nil
}

// FindAll selects all the documents satisfying q.
func (db *DB) FindAll(q *query.Query) ([]*document.Document, error) {
	q, err := normalizeCriteria(q)
	if err != nil {
		return nil, err
	}

	docs := make([]*document.Document, 0)
	err = db.IterateDocs(q, func(doc *document.Document) error {
		docs = append(docs, doc)
		return nil
	})
	return docs, err
}

func (db *DB) IterateDocs(q *query.Query, consumer docConsumer) error {
	tx, err := db.store.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	return db.iterateDocs(tx, q, consumer)
}

// FindFirst returns the first document (if any) satisfying the query.
func (db *DB) FindFirst(q *query.Query) (*document.Document, error) {
	docs, err := db.FindAll(q.Limit(1))

	var doc *document.Document
	if len(docs) > 0 {
		doc = docs[0]
	}
	return doc, err
}

// ForEach runs the consumer function for each document matching the provided query.
// If false is returned from the consumer function, then the iteration is stoppedocument.
func (db *DB) ForEach(q *query.Query, consumer func(_ *document.Document) bool) error {
	q, err := normalizeCriteria(q)
	if err != nil {
		return err
	}

	return db.IterateDocs(q, func(doc *document.Document) error {
		if !consumer(doc) {
			return internal.ErrStopIteration
		}
		return nil
	})
}

// Count returns the number of documents which satisfy the query (i.e. len(q.FindAll()) == q.Count()).
func (db *DB) Count(q *query.Query) (int, error) {
	q, err := normalizeCriteria(q)
	if err != nil {
		return -1, err
	}

	if q.Criteria() == nil { // simply return the size of the collection in this case
		return db.countCollection(q)
	}

	num := 0
	err = db.IterateDocs(q, func(doc *document.Document) error {
		num++
		return nil
	})
	return num, err
}

func (db *DB) countCollection(q *query.Query) (int, error) {
	size, err := db.getCollectionSize(q.Collection())
	size -= q.GetSkip()

	if size < 0 {
		size = 0
	}

	if q.GetLimit() >= 0 && q.GetLimit() < size {
		return q.GetLimit(), err
	}

	return size, err
}

func (db *DB) getCollectionSize(collection string) (int, error) {
	tx, err := db.store.Begin(false)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	meta, err := db.getCollectionMeta(collection, tx)
	if err != nil {
		return -1, err
	}
	return meta.Size, nil
}

// Exists returns true if and only if the query result set is not empty.
func (db *DB) Exists(q *query.Query) (bool, error) {
	doc, err := db.FindFirst(q)
	return doc != nil, err
}

// FindById returns the document with the given id, if such a document exists and satisfies the underlying query, or null.
func (db *DB) FindById(collection string, id string) (*document.Document, error) {
	tx, err := db.store.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ok, err := db.hasCollection(collection, tx)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrCollectionNotExist
	}

	return getDocumentById(collection, id, tx)
}

func getDocumentById(collectionName string, id string, tx store.Tx) (*document.Document, error) {
	value, err := tx.Get([]byte(getDocumentKey(collectionName, id)))
	if value == nil || err != nil {
		return nil, err
	}
	return document.Decode(value)
}

// DeleteById removes the document with the given id from the underlying collection, provided that such a document exists and satisfies the underlying query.
func (db *DB) DeleteById(collection string, id string) error {
	tx, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	meta, err := db.getCollectionMeta(collection, tx)
	if err != nil {
		return err
	}

	indexes := db.getIndexes(tx, collection, meta)

	if err := db.getDocAndDeleteFromIndexes(tx, indexes, collection, id); err != nil {
		return err
	}

	if err := tx.Delete([]byte(getDocumentKey(collection, id))); err != nil {
		return err
	}

	meta.Size--
	if err := db.saveCollectionMetadata(collection, meta, tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) getDocAndDeleteFromIndexes(tx store.Tx, indexes []index.Index, collection string, docId string) error {
	if len(indexes) == 0 {
		return nil
	}

	doc, err := getDocumentById(collection, docId, tx)
	if err != nil {
		return err
	}

	if doc == nil {
		return nil
	}

	for _, idx := range indexes {
		value := doc.Get(idx.Field())
		if err := idx.Remove(doc.ObjectId(), value); err != nil {
			return err
		}
	}
	return nil
}

// UpdateById updates the document with the specified id using the supplied update map.
// If no document with the specified id exists, an ErrDocumentNotExist is returnedocument.
func (db *DB) UpdateById(collectionName string, docId string, updater func(doc *document.Document) *document.Document) error {
	tx, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	meta, err := db.getCollectionMeta(collectionName, tx)
	if err != nil {
		return err
	}

	indexes := db.getIndexes(tx, collectionName, meta)

	docKey := getDocumentKey(collectionName, docId)
	value, err := tx.Get([]byte(docKey))
	if err != nil {
		return err
	}

	if value == nil {
		return ErrDocumentNotExist
	}

	doc, err := document.Decode(value)
	if err != nil {
		return err
	}

	updatedDoc := updater(doc)
	if err := db.updateIndexesOnDocUpdate(tx, indexes, doc, updatedDoc); err != nil {
		return err
	}

	if err := saveDocument(updatedDoc, []byte(docKey), tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) updateIndexesOnDocUpdate(tx store.Tx, indexes []index.Index, oldDoc, newDoc *document.Document) error {
	if err := db.deleteDocFromIndexes(indexes, oldDoc); err != nil {
		return err
	}

	if newDoc != nil {
		if err := db.addDocToIndexes(tx, indexes, newDoc); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) deleteDocFromIndexes(indexes []index.Index, doc *document.Document) error {
	for _, idx := range indexes {
		value := doc.Get(idx.Field())
		if err := idx.Remove(doc.ObjectId(), value); err != nil {
			return err
		}
	}
	return nil
}

// ReplaceById replaces the document with the specified id with the one providedocument.
// If no document exists, an ErrDocumentNotExist is returnedocument.
func (db *DB) ReplaceById(collection, docId string, doc *document.Document) error {
	if doc.ObjectId() != docId {
		return fmt.Errorf("the id of the document must match the one supplied")
	}
	return db.UpdateById(collection, docId, func(_ *document.Document) *document.Document {
		return doc
	})
}

// Update updates all the document selected by q using the provided updateMap.
// Each update is specified by a mapping fieldName -> newValue.
func (db *DB) Update(q *query.Query, updateMap map[string]interface{}) error {
	q, err := normalizeCriteria(q)
	if err != nil {
		return err
	}

	return db.UpdateFunc(q, func(doc *document.Document) *document.Document {
		newDoc := doc.Copy()
		newDoc.SetAll(updateMap)
		return newDoc
	})
}

// UpdateFunc updates all the document selected by q using the provided function.
func (db *DB) UpdateFunc(q *query.Query, updateFunc func(doc *document.Document) *document.Document) error {
	txn, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	q, err = normalizeCriteria(q)
	if err != nil {
		return err
	}
	if err := db.replaceDocs(txn, q, updateFunc); err != nil {
		return err
	}
	return txn.Commit()
}

type docUpdater func(doc *document.Document) *document.Document

func (db *DB) replaceDocs(tx store.Tx, q *query.Query, updater docUpdater) error {
	meta, err := db.getCollectionMeta(q.Collection(), tx)
	if err != nil {
		return err
	}

	indexes := db.getIndexes(tx, q.Collection(), meta)

	deletedDocs := 0
	err = db.iterateDocs(tx, q, func(doc *document.Document) error {
		docKey := []byte(getDocumentKey(q.Collection(), doc.ObjectId()))
		newDoc := updater(doc)

		if err := db.updateIndexesOnDocUpdate(tx, indexes, doc, newDoc); err != nil {
			return err
		}

		if newDoc == nil {
			deletedDocs++
			return tx.Delete(docKey)
		}

		return saveDocument(newDoc, docKey, tx)
	})

	if err != nil {
		return err
	}

	if deletedDocs > 0 {
		meta.Size -= deletedDocs
		if err := db.saveCollectionMetadata(q.Collection(), meta, tx); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) iterateDocs(tx store.Tx, q *query.Query, consumer docConsumer) error {
	meta, err := db.getCollectionMeta(q.Collection(), tx)
	if err != nil {
		return err
	}
	nd := buildQueryPlan(q, db.getIndexes(tx, q.Collection(), meta), &consumerNode{consumer: consumer})
	return execPlan(nd, tx)
}

// Delete removes all the documents selected by q from the underlying collection.
func (db *DB) Delete(q *query.Query) error {
	q, err := normalizeCriteria(q)
	if err != nil {
		return err
	}

	tx, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := db.replaceDocs(tx, q, func(_ *document.Document) *document.Document { return nil }); err != nil {
		return err
	}
	return tx.Commit()
}

// ListCollections returns a slice of strings containing the name of each collection stored in the db.
func (db *DB) ListCollections() ([]string, error) {
	tx, err := db.store.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	collections := make([]string, 0)

	prefix := []byte(getCollectionKeyPrefix())
	err = iteratePrefix(prefix, tx, func(item store.Item) error {
		collectionName := string(bytes.TrimPrefix(item.Key, prefix))
		collections = append(collections, collectionName)
		return nil
	})
	return collections, err
}

func iteratePrefix(prefix []byte, tx store.Tx, itemConsumer func(item store.Item) error) error {
	cursor, err := tx.Cursor(true)
	if err != nil {
		return err
	}
	defer cursor.Close()

	if err := cursor.Seek(prefix); err != nil {
		return err
	}

	for ; cursor.Valid(); cursor.Next() {
		item, err := cursor.Item()
		if err != nil {
			return err
		}

		if !bytes.HasPrefix(item.Key, prefix) {
			return nil
		}
		err = itemConsumer(item)

		// do not propagate iteration stop error
		if errors.Is(err, internal.ErrStopIteration) {
			return nil
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// CreateIndex creates an index for the specified for the specified (index, collection) pair.
func (db *DB) CreateIndex(collection, field string) error {
	return db.createIndex(collection, field, index.SingleField)
}

func (db *DB) createIndex(collection, field string, indexType index.Type) error {
	tx, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	meta, err := db.getCollectionMeta(collection, tx)
	if err != nil {
		return err
	}

	for i := 0; i < len(meta.Indexes); i++ {
		if meta.Indexes[i].Field == field {
			return ErrIndexExist
		}
	}

	if meta.Indexes == nil {
		meta.Indexes = make([]index.Info, 0)
	}
	meta.Indexes = append(meta.Indexes, index.Info{Field: field, Type: indexType})

	idx := index.CreateIndex(collection, field, indexType, tx)

	err = db.iterateDocs(tx, query.NewQuery(collection), func(doc *document.Document) error {
		value := doc.Get(field)
		return idx.Add(doc.ObjectId(), value, doc.TTL())
	})

	if err != nil {
		return err
	}

	if err := db.saveCollectionMetadata(collection, meta, tx); err != nil {
		return err
	}

	return tx.Commit()
}

// HasIndex returns true if an index exists for the specified (index, collection) pair.
func (db *DB) HasIndex(collection, field string) (bool, error) {
	tx, err := db.store.Begin(false)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	return db.hasIndex(tx, collection, field)
}

func (db *DB) hasIndex(tx store.Tx, collection, field string) (bool, error) {
	meta, err := db.getCollectionMeta(collection, tx)
	if err == nil {
		for _, idx := range meta.Indexes {
			if idx.Field == field {
				return true, nil
			}
		}
	}
	return false, err
}

// DropIndex deletes the index, is such index exists for the specified (index, collection) pair.
func (db *DB) DropIndex(collection, field string) error {
	txn, err := db.store.Begin(true)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	meta, err := db.getCollectionMeta(collection, txn)
	if err != nil {
		return err
	}

	j := -1
	for i := 0; i < len(meta.Indexes); i++ {
		if meta.Indexes[i].Field == field {
			j = i
		}
	}

	if j < 0 {
		return ErrIndexNotExist
	}

	idxType := meta.Indexes[j].Type

	meta.Indexes[j] = meta.Indexes[0]
	meta.Indexes = meta.Indexes[1:]

	idx := index.CreateIndex(collection, field, idxType, txn)

	if err := idx.Drop(); err != nil {
		return err
	}

	if err := db.saveCollectionMetadata(collection, meta, txn); err != nil {
		return err
	}
	return txn.Commit()
}

// ListIndexes returns a list containing the names of all the indexes for the specified collection.
func (db *DB) ListIndexes(collection string) ([]index.Info, error) {
	txn, err := db.store.Begin(false)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	return db.listIndexes(collection, txn)
}

func (db *DB) listIndexes(collection string, tx store.Tx) ([]index.Info, error) {
	meta, err := db.getCollectionMeta(collection, tx)
	return meta.Indexes, err
}

func normalizeCriteria(q *query.Query) (*query.Query, error) {
	if q.Criteria() != nil {
		v := &CriteriaNormalizeVisitor{}
		c := q.Criteria().Accept(v)

		if v.err != nil {
			return nil, v.err
		}

		q = q.Where(c.(query.Criteria))
	}
	return q, nil
}
