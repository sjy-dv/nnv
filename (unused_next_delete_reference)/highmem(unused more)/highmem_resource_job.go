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

package highmem

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type HighMem struct {
	Collections  map[string]*CollectionMem
	groupLock    sync.RWMutex
	commiter     *time.Ticker
	stopCommiter chan bool
}

type CollectionMem struct {
	Data            map[uint64]interface{}
	CollectionName  string
	Distance        string
	Quantization    string
	Dim             uint32
	Connectivity    uint32
	ExpansionAdd    uint32
	ExpansionSearch uint32
	Multi           bool
	Storage         string
	collectionLock  sync.RWMutex
}

type CollectionConfig struct {
	CollectionName  string
	Distance        string
	Quantization    string
	Dim             uint32
	Connectivity    uint32
	ExpansionAdd    uint32
	ExpansionSearch uint32
	Multi           bool
	Storage         string
}

type CollectionInfo struct {
	CollectionName  string
	Distance        string
	Quantization    string
	Dim             uint32
	Connectivity    uint32
	ExpansionAdd    uint32
	ExpansionSearch uint32
	Multi           bool
	DataSize        int
	Storage         string
}

func NewHighMemory() *HighMem {
	return &HighMem{
		Collections:  map[string]*CollectionMem{},
		stopCommiter: make(chan bool),
	}
}

// return nil not exists collection
func (xx *HighMem) getCollection(collectionName string) *CollectionMem {
	xx.groupLock.RLock()
	col, exists := xx.Collections[collectionName]
	xx.groupLock.RUnlock()
	if exists {
		return col
	}
	return nil
}

func (xx *HighMem) existsCollection(collectionName string) bool {
	xx.groupLock.RLock()
	_, exists := xx.Collections[collectionName]
	xx.groupLock.RUnlock()
	return exists
}

func (xx *HighMem) CreateCollection(collectionName string, cfg CollectionConfig) error {
	c := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				c <- fmt.Errorf(panicr, r)
			}
		}()
		// check col
		ok := xx.existsCollection(collectionName)
		if ok {
			// already collection exists
			c <- errors.New("already exists collection")
			return
		}
		xx.groupLock.Lock()
		xx.Collections[collectionName] = &CollectionMem{
			Data:            make(map[uint64]interface{}),
			CollectionName:  collectionName,
			Distance:        cfg.Distance,
			Quantization:    cfg.Quantization,
			Dim:             cfg.Dim,
			ExpansionAdd:    cfg.ExpansionAdd,
			ExpansionSearch: cfg.ExpansionSearch,
			Multi:           cfg.Multi,
			Storage:         cfg.Storage,
		}
		xx.groupLock.Unlock()

		//=========== vector build============//
		err := tensorLinker.CreateTensorIndex(collectionName, cfg)
		if err != nil {
			xx.groupLock.Lock()
			delete(xx.Collections, collectionName)
			xx.groupLock.Unlock()
			c <- err
			return
		}
		//==========bitmap index build=======//
		err = indexdb.CreateIndex(collectionName)
		if err != nil {
			xx.groupLock.Lock()
			delete(xx.Collections, collectionName)
			xx.groupLock.Unlock()
			c <- tensorLinker.DropTensorIndex(collectionName)
			return
		}
		stateManager.loadchecker.clcLock.Lock()
		stateManager.loadchecker.collections[collectionName] = true
		stateManager.loadchecker.clcLock.Unlock()
		stateManager.auth.authLock.Lock()
		stateManager.auth.collections[collectionName] = true
		stateManager.auth.authLock.Unlock()
		c <- nil
	}()
	return <-c
}

func (xx *HighMem) DropCollection(collectionName string) error {
	c := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				c <- fmt.Errorf(panicr, r)
			}
		}()
		ok := xx.existsCollection(collectionName)
		if !ok {
			c <- nil
			return
		}
		xx.groupLock.Lock()
		delete(xx.Collections, collectionName)
		xx.groupLock.Unlock()
		err := tensorLinker.DropTensorIndex(collectionName)
		if err != nil {
			c <- err
			return
		}
		c <- indexdb.DropIndex(collectionName)
	}()
	return <-c
}

func (xx *HighMem) GetCollection(collectionName string) (CollectionConfig, error) {
	type cc struct {
		Result CollectionConfig
		Error  error
	}
	c := make(chan cc, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				c <- cc{
					Error: fmt.Errorf(panicr, r),
				}
			}
		}()
		col := xx.getCollection(collectionName)
		if col == nil {
			c <- cc{
				Error: fmt.Errorf("not found collection: %s. check load collection first", collectionName),
			}
			return
		}
		c <- cc{
			Result: CollectionConfig{
				CollectionName:  col.CollectionName,
				Distance:        col.Distance,
				Quantization:    col.Quantization,
				Dim:             col.Dim,
				ExpansionAdd:    col.ExpansionAdd,
				ExpansionSearch: col.ExpansionSearch,
				Multi:           col.Multi,
				Storage:         col.Storage,
			},
		}
	}()

	out := <-c
	return out.Result, out.Error
}

func (xx *HighMem) GetCollections() ([]CollectionConfig, error) {
	type cc struct {
		Result []CollectionConfig
		Error  error
	}
	c := make(chan cc, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				c <- cc{
					Error: fmt.Errorf(panicr, r),
				}
			}
		}()
		ccs := make([]CollectionConfig, 0)
		xx.groupLock.RLock()
		for collectionName := range xx.Collections {
			col := xx.getCollection(collectionName)
			if col == nil {
				c <- cc{
					Error: fmt.Errorf("not found collection %s check laod collection first", collectionName),
				}
				return
			}
			ccs = append(ccs, CollectionConfig{
				CollectionName:  col.CollectionName,
				Distance:        col.Distance,
				Quantization:    col.Quantization,
				Dim:             col.Dim,
				ExpansionAdd:    col.ExpansionAdd,
				ExpansionSearch: col.ExpansionSearch,
				Multi:           col.Multi,
				Storage:         col.Storage,
			})
		}
		xx.groupLock.RUnlock()
		c <- cc{
			Result: ccs,
		}
	}()
	out := <-c
	return out.Result, out.Error
}

func (xx *HighMem) LoadCollection(collectionName string) (CollectionInfo, error) {
	xx.groupLock.RLock()
	collections, exists := xx.Collections[collectionName]
	xx.groupLock.RUnlock()
	if exists {
		return CollectionInfo{
			CollectionName:  collections.CollectionName,
			Distance:        collections.Distance,
			Quantization:    collections.Distance,
			Dim:             collections.Dim,
			Connectivity:    collections.Connectivity,
			ExpansionAdd:    collections.ExpansionAdd,
			ExpansionSearch: collections.ExpansionSearch,
			Multi:           collections.Multi,
			DataSize:        len(collections.Data),
			Storage:         collections.Storage,
		}, nil
	}
	loadData, err := xx.LoadCommitData(collectionName)
	if err != nil {
		return CollectionInfo{}, err
	}
	loadConfig, err := xx.LoadCommitCollectionConfig(collectionName)
	if err != nil {
		return CollectionInfo{}, err
	}
	// loadindex
	err = xx.LoadCommitIndex(collectionName)
	if err != nil {
		return CollectionInfo{}, err
	}
	// loadtensor
	err = xx.LoadCommitTensor(collectionName, loadConfig, uint(len(loadData)))
	if err != nil {
		return CollectionInfo{}, err
	}
	// not error
	mergeCommit := &CollectionMem{
		Data:            loadData,
		CollectionName:  collectionName,
		Distance:        loadConfig.Distance,
		Quantization:    loadConfig.Quantization,
		Dim:             loadConfig.Dim,
		Connectivity:    loadConfig.Connectivity,
		ExpansionAdd:    loadConfig.ExpansionAdd,
		ExpansionSearch: loadConfig.ExpansionSearch,
		Multi:           loadConfig.Multi,
		Storage:         loadConfig.Storage,
	}
	xx.groupLock.Lock()
	xx.Collections[collectionName] = &CollectionMem{}
	xx.Collections[collectionName] = mergeCommit
	xx.groupLock.Unlock()
	stateManager.loadchecker.clcLock.Lock()
	stateManager.loadchecker.collections[collectionName] = true
	stateManager.loadchecker.clcLock.Unlock()
	stateManager.auth.authLock.Lock()
	stateManager.auth.collections[collectionName] = true
	stateManager.auth.authLock.Unlock()
	return CollectionInfo{
		CollectionName:  collectionName,
		Distance:        loadConfig.Distance,
		Quantization:    loadConfig.Quantization,
		Dim:             loadConfig.Dim,
		Connectivity:    loadConfig.Connectivity,
		ExpansionAdd:    loadConfig.ExpansionAdd,
		ExpansionSearch: loadConfig.ExpansionSearch,
		Multi:           loadConfig.Multi,
		DataSize:        len(mergeCommit.Data),
		Storage:         loadConfig.Storage,
	}, nil
}

func (xx *HighMem) ReleaseCollection(collectionName string) error {
	xx.groupLock.RLock()
	_, exists := xx.Collections[collectionName]
	xx.groupLock.RUnlock()
	if !exists {
		return nil
	}
	//if exists
	err := xx.CommitData(collectionName)
	if err != nil {
		return err
	}
	err = xx.CommitCollectionConfig(collectionName)
	if err != nil {
		return err
	}
	err = xx.CommitIndex(collectionName)
	if err != nil {
		return err
	}
	err = xx.CommitTensor(collectionName)
	if err != nil {
		return err
	}
	//release memory
	xx.groupLock.Lock()
	delete(xx.Collections, collectionName)
	xx.groupLock.Unlock()

	indexdb.indexLock.Lock()
	delete(indexdb.indexes, collectionName)
	indexdb.indexLock.Unlock()

	tensorLinker.tensorLock.Lock()
	//release c++ runtime
	err = tensorLinker.tensors[collectionName].Destroy()
	if err != nil {
		tensorLinker.tensorLock.Unlock()
		return err
	}
	delete(tensorLinker.tensors, collectionName)
	tensorLinker.tensorLock.Unlock()

	stateManager.loadchecker.clcLock.Lock()
	stateManager.loadchecker.collections[collectionName] = false
	stateManager.loadchecker.clcLock.Unlock()
	stateManager.auth.authLock.Lock()
	stateManager.auth.collections[collectionName] = false
	stateManager.auth.authLock.Unlock()
	return nil
}
