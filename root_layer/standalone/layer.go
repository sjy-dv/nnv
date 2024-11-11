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

package standalone

import (
	"context"
	"time"

	"github.com/sjy-dv/nnv/gen/protoc/v2/dataCoordinatorV2"
	"github.com/sjy-dv/nnv/gen/protoc/v2/resourceCoordinatorV2"
	"github.com/sjy-dv/nnv/highmem"
	"google.golang.org/grpc"
)

type RootLayer struct {
	Ctx    context.Context
	Cancel context.CancelFunc

	HighMem *highmem.HighMem
	// StreamLayer    *nats.Conn
	// StreamLayerCtx nats.JetStreamContext

	// VBucket     *hnsw.HnswBucket // vector store
	// BitmapIndex *index.BitmapIndex
	S *grpc.Server
}

type rpcLayer struct {
	X1 *datasetCoordinator
	X2 *resourceCoordinator
	// rootClone *RootLayer
}

type datasetCoordinator struct {
	dataCoordinatorV2.UnimplementedDatasetCoordinatorServer
	rpcLayer
}

type resourceCoordinator struct {
	resourceCoordinatorV2.UnimplementedResourceCoordinatorServer
	rpcLayer
}

const (
	b  = 1
	kb = 1024
	mb = 1024 * 1024
	gb = 1024 * 1024 * 1024

	B  = 1
	KB = 1024
	MB = 1024 * 1024
	GB = 1024 * 1024 * 1024
)

const DefaultMsgSize = 104858000 // 10mb
const DefaultKeepAliveTimeout = 10 * time.Second
const DefaultKeepAlive = 60 * time.Second
const DefaultEnforcementPolicyMinTime = 5 * time.Second

var UncaughtPanicError = "uncaught panic error: %v"
