package backup

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/sjy-dv/nnv/backup/document"
	"github.com/sjy-dv/nnv/backup/query"
)

// ExportCollection exports an existing collection to a JSON file.
func (db *DB) ExportCollection(collectionName string, exportPath string) error {
	exists, err := db.HasCollection(collectionName)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCollectionNotExist
	}

	result, err := db.FindAll(query.NewQuery(collectionName))
	if err != nil {
		return err
	}

	docs := make([]map[string]interface{}, 0)
	for _, doc := range result {
		docs = append(docs, doc.AsMap())
	}

	jsonString, err := json.Marshal(docs)
	if err != nil {
		return err
	}

	return os.WriteFile(exportPath, jsonString, os.ModePerm)
}

// ImportCollection imports a collection from a JSON file.
func (db *DB) ImportCollection(collectionName string, importPath string) error {
	file, err := os.Open(importPath)
	if err != nil {
		return err
	}

	if err := db.CreateCollection(collectionName); err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	jsonObjects := make([]*map[string]interface{}, 0)
	err = json.NewDecoder(reader).Decode(&jsonObjects)
	if err != nil {
		return err
	}

	docs := make([]*document.Document, 0)
	for _, doc := range jsonObjects {
		docs = append(docs, document.NewDocumentOf(*doc))
	}
	return db.Insert(collectionName, docs...)
}