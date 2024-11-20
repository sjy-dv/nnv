package edge

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

type vectorspace interface {
	// CreateCollection(config CollectionConfig) error
	// DropCollection(collectionName string) error
	InsertVector(collectionName string, commitId uint64, vector Vector) error
	UpdateVector(collectionName string, id uint64, vector Vector) error
	RemoveVector(collectionName string, id uint64) error
	FullScan(collectionName string, target Vector, topK int) (*ResultSet, error)
}

type Vectorstore struct {
	Space map[string]vectorspace
	slock sync.RWMutex
}

func NewVectorstore() *Vectorstore {
	return &Vectorstore{
		Space: make(map[string]vectorspace),
	}
}

func (xx *Vectorstore) CreateCollection(config CollectionConfig) error {
	xx.slock.RLock()
	_, ok := xx.Space[config.CollectionName]
	xx.slock.RUnlock()
	if ok {
		return fmt.Errorf(ErrCollectionExists, config.CollectionName)
	}
	var vectorstore vectorspace
	if config.Quantization == F8_QUANTIZATION {
		vectorstore = newF8Vectorstore(config)
	} else if config.Quantization == F16_QUANTIZATION {
		vectorstore = newF16Vectorstore(config)
	} else if config.Quantization == BF16_QUANTIZATION {
		vectorstore = newBF16Vectorstore(config)
	} else if config.Quantization == NONE_QAUNTIZATION {
		vectorstore = newSimpleVectorstore(config)
	} else {
		return errors.New("not support quantization type")
	}
	xx.slock.Lock()
	xx.Space[config.CollectionName] = vectorstore
	xx.slock.Unlock()
	return nil
}

func (xx *Vectorstore) DropCollection(collectionName string) error {
	xx.slock.RLock()
	_, ok := xx.Space[collectionName]
	xx.slock.RUnlock()
	if !ok {
		return nil
	}

	xx.slock.Lock()
	defer xx.slock.Unlock()
	delete(xx.Space, collectionName)
	return nil
}

func (xx *Vectorstore) InsertVector(collectionName string, commitId uint64, vector Vector) error {
	xx.slock.RLock()
	basis, ok := xx.Space[collectionName]
	xx.slock.RUnlock()
	if !ok {
		return fmt.Errorf(ErrCollectionNotFound, collectionName)
	}
	return basis.InsertVector(collectionName, commitId, vector)
}

func (xx *Vectorstore) UpdateVector(collectionName string, id uint64, vector Vector) error {
	xx.slock.RLock()
	basis, ok := xx.Space[collectionName]
	xx.slock.RUnlock()
	if !ok {
		return fmt.Errorf(ErrCollectionNotFound, collectionName)
	}
	return basis.UpdateVector(collectionName, id, vector)
}

func (xx *Vectorstore) RemoveVector(collectionName string, id uint64) error {
	xx.slock.RLock()
	basis, ok := xx.Space[collectionName]
	xx.slock.RUnlock()
	if !ok {
		return fmt.Errorf(ErrCollectionNotFound, collectionName)
	}
	return basis.RemoveVector(collectionName, id)
}

func (xx *Vectorstore) FullScan(collectionName string, target Vector, topK int,
) (*ResultSet, error) {
	xx.slock.RLock()
	basis, ok := xx.Space[collectionName]
	xx.slock.RUnlock()
	if !ok {
		return nil, fmt.Errorf(ErrCollectionNotFound, collectionName)
	}
	return basis.FullScan(collectionName, target, topK)
}

func (xx *Vectorstore) Commit(collectionName string) error {

	return nil
}

func (xx *Vectorstore) Load(collectionName string, config CollectionConfig) error {
	xx.slock.Lock()
	defer xx.slock.Unlock()
	if config.Quantization == F8_QUANTIZATION {
		xx.Space[collectionName] = newF8Vectorstore(config)
	} else if config.Quantization == F16_QUANTIZATION {
		xx.Space[collectionName] = newF16Vectorstore(config)
	} else if config.Quantization == BF16_QUANTIZATION {
		xx.Space[collectionName] = newBF16Vectorstore(config)
	} else if config.Quantization == NONE_QAUNTIZATION {
		xx.Space[collectionName] = newSimpleVectorstore(config)
	} else {
		return errors.New("not support quantization type")
	}
	return nil
}

func Normalize(vec Vector) Vector {
	var norm float32
	out := make([]float32, len(vec))
	for i := range vec {
		norm += vec[i] * vec[i]
	}
	if norm == 0 {
		return out
	}

	norm = float32(math.Sqrt(float64(norm)))
	for i := range vec {
		out[i] = vec[i] / norm
	}

	return out
}
