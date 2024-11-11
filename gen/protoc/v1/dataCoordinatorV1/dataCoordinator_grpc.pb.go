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

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0
// source: idl/proto/v1/dataCoordinator.proto

package dataCoordinatorV1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DatasetCoordinator_Ping_FullMethodName                  = "/dataCoordinatorV1.DatasetCoordinator/Ping"
	DatasetCoordinator_Insert_FullMethodName                = "/dataCoordinatorV1.DatasetCoordinator/Insert"
	DatasetCoordinator_Update_FullMethodName                = "/dataCoordinatorV1.DatasetCoordinator/Update"
	DatasetCoordinator_Delete_FullMethodName                = "/dataCoordinatorV1.DatasetCoordinator/Delete"
	DatasetCoordinator_VectorSearch_FullMethodName          = "/dataCoordinatorV1.DatasetCoordinator/VectorSearch"
	DatasetCoordinator_FilterSearch_FullMethodName          = "/dataCoordinatorV1.DatasetCoordinator/FilterSearch"
	DatasetCoordinator_HybridSearch_FullMethodName          = "/dataCoordinatorV1.DatasetCoordinator/HybridSearch"
	DatasetCoordinator_BatchInsert_FullMethodName           = "/dataCoordinatorV1.DatasetCoordinator/BatchInsert"
	DatasetCoordinator_BatchUpdate_FullMethodName           = "/dataCoordinatorV1.DatasetCoordinator/BatchUpdate"
	DatasetCoordinator_BatchDelete_FullMethodName           = "/dataCoordinatorV1.DatasetCoordinator/BatchDelete"
	DatasetCoordinator_Put_FullMethodName                   = "/dataCoordinatorV1.DatasetCoordinator/Put"
	DatasetCoordinator_Get_FullMethodName                   = "/dataCoordinatorV1.DatasetCoordinator/Get"
	DatasetCoordinator_Iterator_FullMethodName              = "/dataCoordinatorV1.DatasetCoordinator/Iterator"
	DatasetCoordinator_Del_FullMethodName                   = "/dataCoordinatorV1.DatasetCoordinator/Del"
	DatasetCoordinator_PerformanceCompaction_FullMethodName = "/dataCoordinatorV1.DatasetCoordinator/PerformanceCompaction"
)

// DatasetCoordinatorClient is the client API for DatasetCoordinator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// dataset coordinator is managing kv data or vector data
type DatasetCoordinatorClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// vector rpc
	Insert(ctx context.Context, in *ModifyDataset, opts ...grpc.CallOption) (*Response, error)
	Update(ctx context.Context, in *ModifyDataset, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *DeleteDataset, opts ...grpc.CallOption) (*Response, error)
	// vectorsearch <- only search vector query
	// filtersearch <- not using vector, only use filter
	// hybrid filter + vector
	VectorSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error)
	FilterSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error)
	HybridSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error)
	BatchInsert(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StreamModifyDataset, Response], error)
	BatchUpdate(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StreamModifyDataset, Response], error)
	BatchDelete(ctx context.Context, in *BatchDeleteIds, opts ...grpc.CallOption) (*Response, error)
	// kv rpc
	Put(ctx context.Context, in *ModifyKV, opts ...grpc.CallOption) (*Response, error)
	Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*GetValue, error)
	Iterator(ctx context.Context, in *Key, opts ...grpc.CallOption) (*GetValues, error)
	Del(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Response, error)
	// It may conflict with multiple tasks, so we recommend avoiding it during busy times and running it during idle times.
	PerformanceCompaction(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Response, error)
}

type datasetCoordinatorClient struct {
	cc grpc.ClientConnInterface
}

func NewDatasetCoordinatorClient(cc grpc.ClientConnInterface) DatasetCoordinatorClient {
	return &datasetCoordinatorClient{cc}
}

func (c *datasetCoordinatorClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DatasetCoordinator_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) Insert(ctx context.Context, in *ModifyDataset, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, DatasetCoordinator_Insert_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) Update(ctx context.Context, in *ModifyDataset, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, DatasetCoordinator_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) Delete(ctx context.Context, in *DeleteDataset, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, DatasetCoordinator_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) VectorSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, DatasetCoordinator_VectorSearch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) FilterSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, DatasetCoordinator_FilterSearch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) HybridSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, DatasetCoordinator_HybridSearch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) BatchInsert(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StreamModifyDataset, Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DatasetCoordinator_ServiceDesc.Streams[0], DatasetCoordinator_BatchInsert_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StreamModifyDataset, Response]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DatasetCoordinator_BatchInsertClient = grpc.BidiStreamingClient[StreamModifyDataset, Response]

func (c *datasetCoordinatorClient) BatchUpdate(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StreamModifyDataset, Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DatasetCoordinator_ServiceDesc.Streams[1], DatasetCoordinator_BatchUpdate_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StreamModifyDataset, Response]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DatasetCoordinator_BatchUpdateClient = grpc.BidiStreamingClient[StreamModifyDataset, Response]

func (c *datasetCoordinatorClient) BatchDelete(ctx context.Context, in *BatchDeleteIds, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, DatasetCoordinator_BatchDelete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) Put(ctx context.Context, in *ModifyKV, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, DatasetCoordinator_Put_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*GetValue, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetValue)
	err := c.cc.Invoke(ctx, DatasetCoordinator_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) Iterator(ctx context.Context, in *Key, opts ...grpc.CallOption) (*GetValues, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetValues)
	err := c.cc.Invoke(ctx, DatasetCoordinator_Iterator_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) Del(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, DatasetCoordinator_Del_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetCoordinatorClient) PerformanceCompaction(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, DatasetCoordinator_PerformanceCompaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatasetCoordinatorServer is the server API for DatasetCoordinator service.
// All implementations should embed UnimplementedDatasetCoordinatorServer
// for forward compatibility.
//
// dataset coordinator is managing kv data or vector data
type DatasetCoordinatorServer interface {
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// vector rpc
	Insert(context.Context, *ModifyDataset) (*Response, error)
	Update(context.Context, *ModifyDataset) (*Response, error)
	Delete(context.Context, *DeleteDataset) (*Response, error)
	// vectorsearch <- only search vector query
	// filtersearch <- not using vector, only use filter
	// hybrid filter + vector
	VectorSearch(context.Context, *SearchReq) (*SearchResponse, error)
	FilterSearch(context.Context, *SearchReq) (*SearchResponse, error)
	HybridSearch(context.Context, *SearchReq) (*SearchResponse, error)
	BatchInsert(grpc.BidiStreamingServer[StreamModifyDataset, Response]) error
	BatchUpdate(grpc.BidiStreamingServer[StreamModifyDataset, Response]) error
	BatchDelete(context.Context, *BatchDeleteIds) (*Response, error)
	// kv rpc
	Put(context.Context, *ModifyKV) (*Response, error)
	Get(context.Context, *Key) (*GetValue, error)
	Iterator(context.Context, *Key) (*GetValues, error)
	Del(context.Context, *Key) (*Response, error)
	// It may conflict with multiple tasks, so we recommend avoiding it during busy times and running it during idle times.
	PerformanceCompaction(context.Context, *emptypb.Empty) (*Response, error)
}

// UnimplementedDatasetCoordinatorServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDatasetCoordinatorServer struct{}

func (UnimplementedDatasetCoordinatorServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDatasetCoordinatorServer) Insert(context.Context, *ModifyDataset) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (UnimplementedDatasetCoordinatorServer) Update(context.Context, *ModifyDataset) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDatasetCoordinatorServer) Delete(context.Context, *DeleteDataset) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDatasetCoordinatorServer) VectorSearch(context.Context, *SearchReq) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VectorSearch not implemented")
}
func (UnimplementedDatasetCoordinatorServer) FilterSearch(context.Context, *SearchReq) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterSearch not implemented")
}
func (UnimplementedDatasetCoordinatorServer) HybridSearch(context.Context, *SearchReq) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HybridSearch not implemented")
}
func (UnimplementedDatasetCoordinatorServer) BatchInsert(grpc.BidiStreamingServer[StreamModifyDataset, Response]) error {
	return status.Errorf(codes.Unimplemented, "method BatchInsert not implemented")
}
func (UnimplementedDatasetCoordinatorServer) BatchUpdate(grpc.BidiStreamingServer[StreamModifyDataset, Response]) error {
	return status.Errorf(codes.Unimplemented, "method BatchUpdate not implemented")
}
func (UnimplementedDatasetCoordinatorServer) BatchDelete(context.Context, *BatchDeleteIds) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDelete not implemented")
}
func (UnimplementedDatasetCoordinatorServer) Put(context.Context, *ModifyKV) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedDatasetCoordinatorServer) Get(context.Context, *Key) (*GetValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedDatasetCoordinatorServer) Iterator(context.Context, *Key) (*GetValues, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Iterator not implemented")
}
func (UnimplementedDatasetCoordinatorServer) Del(context.Context, *Key) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedDatasetCoordinatorServer) PerformanceCompaction(context.Context, *emptypb.Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PerformanceCompaction not implemented")
}
func (UnimplementedDatasetCoordinatorServer) testEmbeddedByValue() {}

// UnsafeDatasetCoordinatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatasetCoordinatorServer will
// result in compilation errors.
type UnsafeDatasetCoordinatorServer interface {
	mustEmbedUnimplementedDatasetCoordinatorServer()
}

func RegisterDatasetCoordinatorServer(s grpc.ServiceRegistrar, srv DatasetCoordinatorServer) {
	// If the following call pancis, it indicates UnimplementedDatasetCoordinatorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DatasetCoordinator_ServiceDesc, srv)
}

func _DatasetCoordinator_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyDataset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_Insert_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).Insert(ctx, req.(*ModifyDataset))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyDataset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).Update(ctx, req.(*ModifyDataset))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDataset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).Delete(ctx, req.(*DeleteDataset))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_VectorSearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).VectorSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_VectorSearch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).VectorSearch(ctx, req.(*SearchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_FilterSearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).FilterSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_FilterSearch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).FilterSearch(ctx, req.(*SearchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_HybridSearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).HybridSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_HybridSearch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).HybridSearch(ctx, req.(*SearchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_BatchInsert_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DatasetCoordinatorServer).BatchInsert(&grpc.GenericServerStream[StreamModifyDataset, Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DatasetCoordinator_BatchInsertServer = grpc.BidiStreamingServer[StreamModifyDataset, Response]

func _DatasetCoordinator_BatchUpdate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DatasetCoordinatorServer).BatchUpdate(&grpc.GenericServerStream[StreamModifyDataset, Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DatasetCoordinator_BatchUpdateServer = grpc.BidiStreamingServer[StreamModifyDataset, Response]

func _DatasetCoordinator_BatchDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchDeleteIds)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).BatchDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_BatchDelete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).BatchDelete(ctx, req.(*BatchDeleteIds))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyKV)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_Put_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).Put(ctx, req.(*ModifyKV))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).Get(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_Iterator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).Iterator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_Iterator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).Iterator(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_Del_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).Del(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetCoordinator_PerformanceCompaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetCoordinatorServer).PerformanceCompaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatasetCoordinator_PerformanceCompaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetCoordinatorServer).PerformanceCompaction(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// DatasetCoordinator_ServiceDesc is the grpc.ServiceDesc for DatasetCoordinator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DatasetCoordinator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dataCoordinatorV1.DatasetCoordinator",
	HandlerType: (*DatasetCoordinatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _DatasetCoordinator_Ping_Handler,
		},
		{
			MethodName: "Insert",
			Handler:    _DatasetCoordinator_Insert_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _DatasetCoordinator_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DatasetCoordinator_Delete_Handler,
		},
		{
			MethodName: "VectorSearch",
			Handler:    _DatasetCoordinator_VectorSearch_Handler,
		},
		{
			MethodName: "FilterSearch",
			Handler:    _DatasetCoordinator_FilterSearch_Handler,
		},
		{
			MethodName: "HybridSearch",
			Handler:    _DatasetCoordinator_HybridSearch_Handler,
		},
		{
			MethodName: "BatchDelete",
			Handler:    _DatasetCoordinator_BatchDelete_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _DatasetCoordinator_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DatasetCoordinator_Get_Handler,
		},
		{
			MethodName: "Iterator",
			Handler:    _DatasetCoordinator_Iterator_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _DatasetCoordinator_Del_Handler,
		},
		{
			MethodName: "PerformanceCompaction",
			Handler:    _DatasetCoordinator_PerformanceCompaction_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BatchInsert",
			Handler:       _DatasetCoordinator_BatchInsert_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "BatchUpdate",
			Handler:       _DatasetCoordinator_BatchUpdate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "idl/proto/v1/dataCoordinator.proto",
}
