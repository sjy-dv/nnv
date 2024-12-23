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

package gomath

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"

	"github.com/viterin/vek/vek32"
)

const VECTOR_COMPONENT_BYTES_SIZE = 4

type Vector []float32

func (v Vector) Len() int { return len(v) }

func (v Vector) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Vector) Less(i, j int) bool { return v[i] < v[j] }

func (v Vector) Sort() Vector {
	sort.Sort(v)
	return v
}

func (v Vector) Save(w io.Writer) error {
	for _, val := range v {
		if err := binary.Write(w, binary.BigEndian, val); err != nil {
			return err
		}
	}
	return nil
}

func (v Vector) Load(r io.Reader) error {
	for i := 0; i < len(v); i++ {
		if err := binary.Read(r, binary.BigEndian, &v[i]); err != nil {
			return err
		}
	}
	return nil
}

func (v Vector) Normalize() {
	factor := vek32.Norm(v)
	vek32.DivNumber_Inplace(v, factor)
}

func assertSameDim(i, j *Vector) {
	if len(*i) != len(*j) {
		panic("Vector sizes do not match.")
	}
}

func VectorFromBytes(bytesSlice [][]byte) (Vector, error) {
	vector := make(Vector, len(bytesSlice))
	for i, b := range bytesSlice {
		var element float32
		buf := bytes.NewReader(b)
		err := binary.Read(buf, binary.BigEndian, &element)
		if err != nil {
			return nil, err
		}
		vector[i] = float32(element)
	}
	return vector, nil
}

func ZerosVector(size int) Vector {
	return make(Vector, size)
}

func OnesVector(size int) Vector {
	vector := make(Vector, size)
	for i := 0; i < size; i++ {
		vector[i] = 1
	}
	return vector
}

func Dot(a, b Vector) float32 {
	var dot float32
	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
	}
	return dot
}

func Length(a Vector) float32 {
	return Sqrt(Dot(a, a))
}

func VectorAdd(a, b Vector) Vector {
	assertSameDim(&a, &b)

	result := make(Vector, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] + b[i]
	}
	return result
}

func VectorSubtract(a, b Vector) Vector {
	assertSameDim(&a, &b)

	result := make(Vector, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] - b[i]
	}
	return result
}

func VectorMultiply(a, b Vector) Vector {
	assertSameDim(&a, &b)

	result := make(Vector, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] * b[i]
	}
	return result
}

func VectorDivide(a, b Vector) Vector {
	assertSameDim(&a, &b)

	result := make(Vector, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] / b[i]
	}
	return result
}

func VectorScalarAdd(a Vector, b float32) Vector {
	result := make(Vector, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] + b
	}
	return result
}

func VectorScalarSubtract(a Vector, b float32) Vector {
	result := make(Vector, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] - b
	}
	return result
}

func VectorScalarMultiply(a Vector, b float32) Vector {
	result := make(Vector, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] * b
	}
	return result
}

func VectorScalarDivide(a Vector, b float32) Vector {
	result := make(Vector, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] / b
	}
	return result
}
