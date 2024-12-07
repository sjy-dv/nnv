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
	"compress/flate"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	// "github.com/sjy-dv/nnv/pkg/fasthnsw"
	"github.com/sjy-dv/nnv/pkg/index"
)

func init() {
	gob.Register(uuid.UUID{})
	gob.Register(map[uint64]interface{}{})
}

func (xx *HighMem) CommitAll() {
	log.Info().Msg("starting commit..")
	log.Info().Msg("attempting to commit collection list")
	err := xx.CommitCollection()
	if err != nil {
		log.Error().Err(err).Msg("failed to commit collection list")
	}
	//Efficiently stores only collections with predictable changes that are loaded into memory
	attemptCount := 0
	failedCount := 0
	for collection, ok := range stateManager.loadchecker.collections {
		// only commit true collection
		//Because collections that are false are committed during the Release process.
		if ok {
			attemptCount++
			err := xx.ReleaseCollection(collection)
			if err != nil {
				failedCount++
				// add to commit logger fatal init
				log.Error().Err(err).Msgf("failed to commit collection: %s", collection)
				cloneCollection := xx.getCollection(collection)
				if cloneCollection == nil {
					log.Warn().Msgf("collection: %s is alreay release memory", collection)
					continue
				}
				err := commitLogger.CommitCrashCollection(CollectionConfig{
					CollectionName:  cloneCollection.CollectionName,
					Distance:        cloneCollection.Distance,
					Quantization:    cloneCollection.Quantization,
					Dim:             cloneCollection.Dim,
					Connectivity:    cloneCollection.Connectivity,
					ExpansionAdd:    cloneCollection.ExpansionAdd,
					ExpansionSearch: cloneCollection.ExpansionSearch,
					Multi:           cloneCollection.Multi,
					Storage:         cloneCollection.Storage,
				})
				if err != nil {
					log.Error().Err(err).Msgf("recovery collection config data also broken: %s", collection)
				}
				log.Info().Msgf("recovery data saved collection: %s", collection)
			}
		}
	}
	log.Info().Msgf("all collection commit attempt:%d failed:%d", attemptCount, failedCount)
}

/*
find origin file
if (backup data exist)

	-> remove {origin}-backup.cdat

check (origin) <- if exists

	-> rename {origin}-backup.cdat

new file write {origin}.cdat
*/
func (xx *HighMem) CommitData(collectionName string) error {
	// check origin file
	_, err := os.Stat(fmt.Sprintf(fLinkCdat, collectionName))
	if err != nil {
		// when os is not exist => origin.cdat ok
		//!os.isnotexist => pure error
		if !os.IsNotExist(err) {
			return err
		}
		// this case is origin not exist
		// we guess to also origin-backup not exist
	} else {
		// case origin exist
		// check backup => ok => delete backup
		// origin rename to backup
		_, err := os.Stat(fmt.Sprintf(backupfLinkCdat, collectionName))
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			//case exist backup
			os.Remove(fmt.Sprintf(backupfLinkCdat, collectionName))
			//
		}
		os.Rename(fmt.Sprintf(fLinkCdat, collectionName),
			fmt.Sprintf(backupfLinkCdat, collectionName))
	}
	// create origin
	var iow io.Writer

	backupCdat := make(map[uint64]interface{})
	xx.Collections[collectionName].collectionLock.RLock()
	backupCdat = xx.Collections[collectionName].Data
	xx.Collections[collectionName].collectionLock.RUnlock()

	f, err := os.OpenFile(fmt.Sprintf(fLinkCdat, collectionName), os.O_TRUNC|
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	iow, _ = flate.NewWriter(f, flate.BestCompression)

	enc := gob.NewEncoder(iow)
	if err := enc.Encode(backupCdat); err != nil {
		return err
	}
	if flusher, ok := iow.(interface{ Flush() error }); ok {
		if err := flusher.Flush(); err != nil {
			return err
		}
	}
	if err := iow.(io.Closer).Close(); err != nil {
		return err
	}
	return nil
}

func (xx *HighMem) CommitCollectionConfig(collectionName string) error {

	_, err := os.Stat(fmt.Sprintf(confJson, collectionName))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		_, err := os.Stat(fmt.Sprintf(backupConfJson, collectionName))
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			os.Remove(fmt.Sprintf(backupConfJson, collectionName))
		}
		os.Rename(fmt.Sprintf(confJson, collectionName),
			fmt.Sprintf(backupConfJson, collectionName))
	}

	conf := CollectionConfig{}
	conf.CollectionName = collectionName
	conf.Connectivity = xx.Collections[collectionName].Connectivity
	conf.Dim = xx.Collections[collectionName].Dim
	conf.Distance = xx.Collections[collectionName].Distance
	conf.ExpansionAdd = xx.Collections[collectionName].ExpansionAdd
	conf.ExpansionSearch = xx.Collections[collectionName].ExpansionSearch
	conf.Multi = xx.Collections[collectionName].Multi
	conf.Quantization = xx.Collections[collectionName].Quantization
	conf.Storage = xx.Collections[collectionName].Storage
	f, err := os.Create(fmt.Sprintf(confJson, collectionName))
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.Encode(conf)
	defer f.Close()
	return nil
}

func (xx *HighMem) CommitIndex(collectionName string) error {

	_, err := os.Stat(fmt.Sprintf(indexBin, collectionName))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		_, err := os.Stat(fmt.Sprintf(backupIndexBin, collectionName))
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			os.Remove(fmt.Sprintf(backupIndexBin, collectionName))
		}
		os.Rename(fmt.Sprintf(indexBin, collectionName),
			fmt.Sprintf(backupIndexBin, collectionName))
	}
	return indexdb.indexes[collectionName].SerializeBinary(fmt.Sprintf(indexBin, collectionName))
}

func (xx *HighMem) CommitTensor(collectionName string) error {

	_, err := os.Stat(fmt.Sprintf(tensorLink, collectionName))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		_, err := os.Stat(fmt.Sprintf(backupTensorLink, collectionName))
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			os.Remove(fmt.Sprintf(backupTensorLink, collectionName))
		}
		os.Rename(fmt.Sprintf(tensorLink, collectionName),
			fmt.Sprintf(backupTensorLink, collectionName))
	}
	return tensorLinker.tensors[collectionName].Save(fmt.Sprintf(tensorLink, collectionName))
}

func (xx *HighMem) CommitCollection() error {
	_, err := os.Stat(collectionJson)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		_, err := os.Stat(backupCollectionJson)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			os.Remove(backupCollectionJson)
		}
		os.Rename(collectionJson, backupCollectionJson)
	}
	// collection list [a,b,c]
	// getcollections is only view loadcollection
	// showcollection return collection names
	type collectionJsonF struct {
		Collections []string
	}
	c := collectionJsonF{}
	cols := make([]string, 0)
	for c, ok := range stateManager.checker.collections {
		if ok {
			cols = append(cols, c)
		}
	}
	c.Collections = cols
	f, err := os.Create(collectionJson)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.Encode(c)
	defer f.Close()
	return nil
}

/*
   laod origin collection
   when user request
*/

func (xx *HighMem) LoadCommitCollection() error {

	collectionJsonData, err := os.ReadFile(collectionJson)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	type collectionJsonF struct {
		Collections []string
	}
	c := collectionJsonF{}
	if err := json.Unmarshal(collectionJsonData, &c); err != nil {
		return err
	}
	for _, col := range c.Collections {
		stateManager.checker.collections[col] = true
	}
	return nil
}

func (xx *HighMem) LoadCommitData(collectionName string) (map[uint64]interface{}, error) {
	_, err := os.Stat(fmt.Sprintf(fLinkCdat, collectionName))
	if err != nil {
		if os.IsNotExist(err) {
			// for _, col := range collections {
			// 	if col == collectionName {
			// 		goto EmptyData
			// 	}
			// }
			if stateManager.checker.collections[collectionName] {
				goto EmptyData
			}
			return nil, fmt.Errorf("collection: %s is not defined [Not Found Collection Error]", collectionName)
		}
		return nil, err
	}
	goto ExistsData
EmptyData:
	return make(map[uint64]interface{}), nil
ExistsData:
	commitCdat, err := os.OpenFile(fmt.Sprintf(fLinkCdat, collectionName), os.O_RDONLY, 0777)
	if err != nil {
		// cdat is damaged
		// after add recovery logic
		return nil, err
	}
	cdat := make(map[uint64]interface{})

	var readIo io.Reader

	readIo = flate.NewReader(commitCdat)

	dataDec := gob.NewDecoder(readIo)
	err = dataDec.Decode(&cdat)
	if err != nil {
		// also cdat is damaged guess
		return nil, err
	}
	err = readIo.(io.Closer).Close()
	if err != nil {
		return nil, err
	}
	return cdat, nil
}

func (xx *HighMem) LoadCommitCollectionConfig(collectionName string) (
	CollectionConfig, error) {
	configJsonData, err := os.ReadFile(fmt.Sprintf(confJson, collectionName))
	if err != nil {
		if os.IsNotExist(err) {
			if stateManager.checker.collections[collectionName] {
				//damaged data file
				// recovery logic add
				return CollectionConfig{}, fmt.Errorf("file %s.conf is damaged", collectionName)
			}

			return CollectionConfig{}, fmt.Errorf("collection Config: %s is not defined [Not Found Collection Error]", collectionName)
		}
		return CollectionConfig{}, err
	}
	cfg := CollectionConfig{}
	if err := json.Unmarshal(configJsonData, &cfg); err != nil {
		return CollectionConfig{}, err
	}
	return cfg, nil
}

func (xx *HighMem) LoadCommitIndex(collectionName string) error {
	_, err := os.Stat(fmt.Sprintf(indexBin, collectionName))
	if err != nil {
		if os.IsNotExist(err) {
			if stateManager.checker.collections[collectionName] {
				goto EmptyIndex
			}

			return fmt.Errorf("collection BitmapIndex: %s is not defined [Not Found Collection Error]", collectionName)
		}
		return err
	}
	goto ExistsIndex
EmptyIndex:
	indexdb.indexLock.Lock()
	indexdb.indexes[collectionName] = index.NewBitmapIndex()
	defer indexdb.indexLock.Unlock()
	return nil
ExistsIndex:
	recoveryIndex := index.NewBitmapIndex()
	err = recoveryIndex.DeserializeBinary(fmt.Sprintf(indexBin, collectionName))
	if err != nil {
		// guess damaged file
		return err
	}
	err = recoveryIndex.ValidateIndex()
	if err != nil {
		// validation failed also damaged
		return err
	}
	indexdb.indexLock.Lock()
	indexdb.indexes[collectionName] = recoveryIndex
	defer indexdb.indexLock.Unlock()
	return nil
}

func (xx *HighMem) LoadCommitTensor(collectionName string, cfg CollectionConfig, cap uint) error {
	tcfg := fasthnsw.DefaultConfig(uint(cfg.Dim))
	tcfg.Connectivity = uint(cfg.Connectivity)
	tcfg.ExpansionAdd = uint(cfg.ExpansionAdd)
	tcfg.ExpansionSearch = uint(cfg.ExpansionSearch)
	tcfg.Multi = cfg.Multi
	if cfg.Quantization != "None" {
		tcfg.Quantization = func() fasthnsw.Quantization {
			switch cfg.Quantization {
			case "BF16":
				return fasthnsw.BF16
			case "F16":
				return fasthnsw.F16
			case "F32":
				return fasthnsw.F32
			case "F64":
				return fasthnsw.F64
			case "I8":
				return fasthnsw.I8
			case "B1":
				return fasthnsw.B1
			}
			return fasthnsw.F16
		}()
	}
	tcfg.Metric = func() fasthnsw.Metric {
		switch cfg.Distance {
		case "InnerProduct":
			return fasthnsw.InnerProduct
		case "Cosine":
			return fasthnsw.Cosine
		case "Haversine":
			return fasthnsw.Haversine
		case "Divergence":
			return fasthnsw.Divergence
		case "Pearson":
			return fasthnsw.Pearson
		case "Hamming":
			return fasthnsw.Hamming
		case "Tanimoto":
			return fasthnsw.Tanimoto
		case "Sorensen":
			return fasthnsw.Sorensen
		}
		return fasthnsw.Cosine
	}()
	newTensor, err := fasthnsw.NewIndex(tcfg)
	if err != nil {
		return err
	}
	_, err = os.Stat(fmt.Sprintf(tensorLink, collectionName))
	if err != nil {
		if os.IsNotExist(err) {
			if stateManager.checker.collections[collectionName] {
				goto EmptyTensor
			}
			return fmt.Errorf("collection TensorSearch: %s is not defined [Not Found Collection Error]", collectionName)
		}
		return err
	}
	goto ExistsTensor
EmptyTensor:
	err = newTensor.Reserve(100_000)
	if err != nil {
		return err
	}
	tensorCapacity = 100_000
	tensorLinker.tensorLock.Lock()
	defer tensorLinker.tensorLock.Unlock()
	tensorLinker.tensors[collectionName] = newTensor
	return nil

ExistsTensor:
	err = newTensor.Load(fmt.Sprintf(tensorLink, collectionName))
	if err != nil {
		return err
	}
	err = newTensor.Reserve(cap + 10_000)
	if err != nil {
		return err
	}
	tensorCapacity = cap + 10_000
	tensorLinker.tensorLock.Lock()
	defer tensorLinker.tensorLock.Unlock()
	tensorLinker.tensors[collectionName] = newTensor
	return nil
}
