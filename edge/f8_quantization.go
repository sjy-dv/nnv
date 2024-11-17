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

package edge

import (
	"github.com/sjy-dv/nnv/pkg/compresshelper"
	"github.com/sjy-dv/nnv/pkg/distance"
	"github.com/sjy-dv/nnv/pkg/gomath"
)

type float8Vec []compresshelper.Float8

var _ Quantization[float8Vec] = Float8Quantization{}

type Float8Quantization struct {
	bufx, bufy gomath.Vector
}

func (q Float8Quantization) Similarity(x, y float8Vec, dist distance.Space) float32 {
	if q.bufx == nil {
		q.bufx = make(gomath.Vector, len(x))
		q.bufy = make(gomath.Vector, len(x))
	}
	for i := range x {
		q.bufx[i] = x[i].Float32()
		q.bufy[i] = y[i].Float32()
	}
	return dist.Distance(q.bufx, q.bufy)
}

func (q Float8Quantization) Lower(v gomath.Vector) (float8Vec, error) {
	out := make(float8Vec, len(v))
	for i, x := range v {
		out[i] = compresshelper.F8Fromfloat32(x)
	}
	return out, nil
}

func (q Float8Quantization) Name() string {
	return "float8"
}

func (q Float8Quantization) LowerSize(dim int) int {
	return 2 * dim
}
