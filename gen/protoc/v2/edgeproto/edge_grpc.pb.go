// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0
// source: idl/proto/v2/edge.proto

package edgeproto

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
	EdgeRpc_Ping_FullMethodName              = "/edgeproto.EdgeRpc/Ping"
	EdgeRpc_CreateCollection_FullMethodName  = "/edgeproto.EdgeRpc/CreateCollection"
	EdgeRpc_DeleteCollection_FullMethodName  = "/edgeproto.EdgeRpc/DeleteCollection"
	EdgeRpc_GetCollection_FullMethodName     = "/edgeproto.EdgeRpc/GetCollection"
	EdgeRpc_LoadCollection_FullMethodName    = "/edgeproto.EdgeRpc/LoadCollection"
	EdgeRpc_ReleaseCollection_FullMethodName = "/edgeproto.EdgeRpc/ReleaseCollection"
	EdgeRpc_Flush_FullMethodName             = "/edgeproto.EdgeRpc/Flush"
	EdgeRpc_Insert_FullMethodName            = "/edgeproto.EdgeRpc/Insert"
	EdgeRpc_Update_FullMethodName            = "/edgeproto.EdgeRpc/Update"
	EdgeRpc_Delete_FullMethodName            = "/edgeproto.EdgeRpc/Delete"
	EdgeRpc_VectorSearch_FullMethodName      = "/edgeproto.EdgeRpc/VectorSearch"
	EdgeRpc_FilterSearch_FullMethodName      = "/edgeproto.EdgeRpc/FilterSearch"
	EdgeRpc_HybridSearch_FullMethodName      = "/edgeproto.EdgeRpc/HybridSearch"
)

// EdgeRpcClient is the client API for EdgeRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EdgeRpcClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateCollection(ctx context.Context, in *Collection, opts ...grpc.CallOption) (*CollectionResponse, error)
	DeleteCollection(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*DeleteCollectionResponse, error)
	GetCollection(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*CollectionDetail, error)
	LoadCollection(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*CollectionDetail, error)
	ReleaseCollection(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*Response, error)
	Flush(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*Response, error)
	Insert(ctx context.Context, in *ModifyDataset, opts ...grpc.CallOption) (*Response, error)
	Update(ctx context.Context, in *ModifyDataset, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *DeleteDataset, opts ...grpc.CallOption) (*Response, error)
	VectorSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error)
	FilterSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error)
	HybridSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error)
}

type edgeRpcClient struct {
	cc grpc.ClientConnInterface
}

func NewEdgeRpcClient(cc grpc.ClientConnInterface) EdgeRpcClient {
	return &edgeRpcClient{cc}
}

func (c *edgeRpcClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, EdgeRpc_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) CreateCollection(ctx context.Context, in *Collection, opts ...grpc.CallOption) (*CollectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionResponse)
	err := c.cc.Invoke(ctx, EdgeRpc_CreateCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) DeleteCollection(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*DeleteCollectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteCollectionResponse)
	err := c.cc.Invoke(ctx, EdgeRpc_DeleteCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) GetCollection(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*CollectionDetail, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionDetail)
	err := c.cc.Invoke(ctx, EdgeRpc_GetCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) LoadCollection(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*CollectionDetail, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionDetail)
	err := c.cc.Invoke(ctx, EdgeRpc_LoadCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) ReleaseCollection(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, EdgeRpc_ReleaseCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) Flush(ctx context.Context, in *CollectionName, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, EdgeRpc_Flush_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) Insert(ctx context.Context, in *ModifyDataset, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, EdgeRpc_Insert_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) Update(ctx context.Context, in *ModifyDataset, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, EdgeRpc_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) Delete(ctx context.Context, in *DeleteDataset, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, EdgeRpc_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) VectorSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, EdgeRpc_VectorSearch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) FilterSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, EdgeRpc_FilterSearch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeRpcClient) HybridSearch(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, EdgeRpc_HybridSearch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EdgeRpcServer is the server API for EdgeRpc service.
// All implementations should embed UnimplementedEdgeRpcServer
// for forward compatibility.
type EdgeRpcServer interface {
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	CreateCollection(context.Context, *Collection) (*CollectionResponse, error)
	DeleteCollection(context.Context, *CollectionName) (*DeleteCollectionResponse, error)
	GetCollection(context.Context, *CollectionName) (*CollectionDetail, error)
	LoadCollection(context.Context, *CollectionName) (*CollectionDetail, error)
	ReleaseCollection(context.Context, *CollectionName) (*Response, error)
	Flush(context.Context, *CollectionName) (*Response, error)
	Insert(context.Context, *ModifyDataset) (*Response, error)
	Update(context.Context, *ModifyDataset) (*Response, error)
	Delete(context.Context, *DeleteDataset) (*Response, error)
	VectorSearch(context.Context, *SearchReq) (*SearchResponse, error)
	FilterSearch(context.Context, *SearchReq) (*SearchResponse, error)
	HybridSearch(context.Context, *SearchReq) (*SearchResponse, error)
}

// UnimplementedEdgeRpcServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEdgeRpcServer struct{}

func (UnimplementedEdgeRpcServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedEdgeRpcServer) CreateCollection(context.Context, *Collection) (*CollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCollection not implemented")
}
func (UnimplementedEdgeRpcServer) DeleteCollection(context.Context, *CollectionName) (*DeleteCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCollection not implemented")
}
func (UnimplementedEdgeRpcServer) GetCollection(context.Context, *CollectionName) (*CollectionDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCollection not implemented")
}
func (UnimplementedEdgeRpcServer) LoadCollection(context.Context, *CollectionName) (*CollectionDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadCollection not implemented")
}
func (UnimplementedEdgeRpcServer) ReleaseCollection(context.Context, *CollectionName) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseCollection not implemented")
}
func (UnimplementedEdgeRpcServer) Flush(context.Context, *CollectionName) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Flush not implemented")
}
func (UnimplementedEdgeRpcServer) Insert(context.Context, *ModifyDataset) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (UnimplementedEdgeRpcServer) Update(context.Context, *ModifyDataset) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedEdgeRpcServer) Delete(context.Context, *DeleteDataset) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedEdgeRpcServer) VectorSearch(context.Context, *SearchReq) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VectorSearch not implemented")
}
func (UnimplementedEdgeRpcServer) FilterSearch(context.Context, *SearchReq) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterSearch not implemented")
}
func (UnimplementedEdgeRpcServer) HybridSearch(context.Context, *SearchReq) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HybridSearch not implemented")
}
func (UnimplementedEdgeRpcServer) testEmbeddedByValue() {}

// UnsafeEdgeRpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EdgeRpcServer will
// result in compilation errors.
type UnsafeEdgeRpcServer interface {
	mustEmbedUnimplementedEdgeRpcServer()
}

func RegisterEdgeRpcServer(s grpc.ServiceRegistrar, srv EdgeRpcServer) {
	// If the following call pancis, it indicates UnimplementedEdgeRpcServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EdgeRpc_ServiceDesc, srv)
}

func _EdgeRpc_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_CreateCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Collection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).CreateCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_CreateCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).CreateCollection(ctx, req.(*Collection))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_DeleteCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).DeleteCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_DeleteCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).DeleteCollection(ctx, req.(*CollectionName))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_GetCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).GetCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_GetCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).GetCollection(ctx, req.(*CollectionName))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_LoadCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).LoadCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_LoadCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).LoadCollection(ctx, req.(*CollectionName))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_ReleaseCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).ReleaseCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_ReleaseCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).ReleaseCollection(ctx, req.(*CollectionName))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_Flush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).Flush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_Flush_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).Flush(ctx, req.(*CollectionName))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyDataset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_Insert_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).Insert(ctx, req.(*ModifyDataset))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyDataset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).Update(ctx, req.(*ModifyDataset))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDataset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).Delete(ctx, req.(*DeleteDataset))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_VectorSearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).VectorSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_VectorSearch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).VectorSearch(ctx, req.(*SearchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_FilterSearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).FilterSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_FilterSearch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).FilterSearch(ctx, req.(*SearchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeRpc_HybridSearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeRpcServer).HybridSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeRpc_HybridSearch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeRpcServer).HybridSearch(ctx, req.(*SearchReq))
	}
	return interceptor(ctx, in, info, handler)
}

// EdgeRpc_ServiceDesc is the grpc.ServiceDesc for EdgeRpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EdgeRpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "edgeproto.EdgeRpc",
	HandlerType: (*EdgeRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _EdgeRpc_Ping_Handler,
		},
		{
			MethodName: "CreateCollection",
			Handler:    _EdgeRpc_CreateCollection_Handler,
		},
		{
			MethodName: "DeleteCollection",
			Handler:    _EdgeRpc_DeleteCollection_Handler,
		},
		{
			MethodName: "GetCollection",
			Handler:    _EdgeRpc_GetCollection_Handler,
		},
		{
			MethodName: "LoadCollection",
			Handler:    _EdgeRpc_LoadCollection_Handler,
		},
		{
			MethodName: "ReleaseCollection",
			Handler:    _EdgeRpc_ReleaseCollection_Handler,
		},
		{
			MethodName: "Flush",
			Handler:    _EdgeRpc_Flush_Handler,
		},
		{
			MethodName: "Insert",
			Handler:    _EdgeRpc_Insert_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _EdgeRpc_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _EdgeRpc_Delete_Handler,
		},
		{
			MethodName: "VectorSearch",
			Handler:    _EdgeRpc_VectorSearch_Handler,
		},
		{
			MethodName: "FilterSearch",
			Handler:    _EdgeRpc_FilterSearch_Handler,
		},
		{
			MethodName: "HybridSearch",
			Handler:    _EdgeRpc_HybridSearch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "idl/proto/v2/edge.proto",
}