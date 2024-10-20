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

package space

import (
	"github.com/sjy-dv/nnv/pkg/gomath"
	"github.com/sjy-dv/nnv/pkg/hnsw/simd/avx"
)

type avxSpaceImpl struct{}

func (avxSpaceImpl) EuclideanDistance(a, b gomath.Vector) float32 {
	return avx.EuclideanDistance(a, b)
}

func (avxSpaceImpl) ManhattanDistance(a, b gomath.Vector) float32 {
	return avx.ManhattanDistance(a, b)
}

func (avxSpaceImpl) CosineDistance(a, b gomath.Vector) float32 {
	return avx.CosineDistance(a, b)
}
