// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: news.proto

package news_proto

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
	NewsService_CreateNewsHandler_FullMethodName = "/data.NewsService/CreateNewsHandler"
	NewsService_ShowNewsHandler_FullMethodName   = "/data.NewsService/ShowNewsHandler"
	NewsService_UpdateNewsHandler_FullMethodName = "/data.NewsService/UpdateNewsHandler"
	NewsService_DeleteNewsHandler_FullMethodName = "/data.NewsService/DeleteNewsHandler"
	NewsService_ListNewsHandler_FullMethodName   = "/data.NewsService/ListNewsHandler"
)

// NewsServiceClient is the client API for NewsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NewsServiceClient interface {
	CreateNewsHandler(ctx context.Context, in *CreateNewsRequest, opts ...grpc.CallOption) (*News, error)
	ShowNewsHandler(ctx context.Context, in *NewsId, opts ...grpc.CallOption) (*News, error)
	UpdateNewsHandler(ctx context.Context, in *UpdateNewsRequest, opts ...grpc.CallOption) (*News, error)
	DeleteNewsHandler(ctx context.Context, in *NewsId, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListNewsHandler(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*NewsList, error)
}

type newsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNewsServiceClient(cc grpc.ClientConnInterface) NewsServiceClient {
	return &newsServiceClient{cc}
}

func (c *newsServiceClient) CreateNewsHandler(ctx context.Context, in *CreateNewsRequest, opts ...grpc.CallOption) (*News, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(News)
	err := c.cc.Invoke(ctx, NewsService_CreateNewsHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) ShowNewsHandler(ctx context.Context, in *NewsId, opts ...grpc.CallOption) (*News, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(News)
	err := c.cc.Invoke(ctx, NewsService_ShowNewsHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) UpdateNewsHandler(ctx context.Context, in *UpdateNewsRequest, opts ...grpc.CallOption) (*News, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(News)
	err := c.cc.Invoke(ctx, NewsService_UpdateNewsHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) DeleteNewsHandler(ctx context.Context, in *NewsId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NewsService_DeleteNewsHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) ListNewsHandler(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*NewsList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NewsList)
	err := c.cc.Invoke(ctx, NewsService_ListNewsHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NewsServiceServer is the server API for NewsService service.
// All implementations must embed UnimplementedNewsServiceServer
// for forward compatibility.
type NewsServiceServer interface {
	CreateNewsHandler(context.Context, *CreateNewsRequest) (*News, error)
	ShowNewsHandler(context.Context, *NewsId) (*News, error)
	UpdateNewsHandler(context.Context, *UpdateNewsRequest) (*News, error)
	DeleteNewsHandler(context.Context, *NewsId) (*emptypb.Empty, error)
	ListNewsHandler(context.Context, *GetAllRequest) (*NewsList, error)
	mustEmbedUnimplementedNewsServiceServer()
}

// UnimplementedNewsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNewsServiceServer struct{}

func (UnimplementedNewsServiceServer) CreateNewsHandler(context.Context, *CreateNewsRequest) (*News, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewsHandler not implemented")
}
func (UnimplementedNewsServiceServer) ShowNewsHandler(context.Context, *NewsId) (*News, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShowNewsHandler not implemented")
}
func (UnimplementedNewsServiceServer) UpdateNewsHandler(context.Context, *UpdateNewsRequest) (*News, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNewsHandler not implemented")
}
func (UnimplementedNewsServiceServer) DeleteNewsHandler(context.Context, *NewsId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNewsHandler not implemented")
}
func (UnimplementedNewsServiceServer) ListNewsHandler(context.Context, *GetAllRequest) (*NewsList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNewsHandler not implemented")
}
func (UnimplementedNewsServiceServer) mustEmbedUnimplementedNewsServiceServer() {}
func (UnimplementedNewsServiceServer) testEmbeddedByValue()                     {}

// UnsafeNewsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NewsServiceServer will
// result in compilation errors.
type UnsafeNewsServiceServer interface {
	mustEmbedUnimplementedNewsServiceServer()
}

func RegisterNewsServiceServer(s grpc.ServiceRegistrar, srv NewsServiceServer) {
	// If the following call pancis, it indicates UnimplementedNewsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NewsService_ServiceDesc, srv)
}

func _NewsService_CreateNewsHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNewsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).CreateNewsHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NewsService_CreateNewsHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).CreateNewsHandler(ctx, req.(*CreateNewsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_ShowNewsHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewsId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).ShowNewsHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NewsService_ShowNewsHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).ShowNewsHandler(ctx, req.(*NewsId))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_UpdateNewsHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNewsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).UpdateNewsHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NewsService_UpdateNewsHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).UpdateNewsHandler(ctx, req.(*UpdateNewsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_DeleteNewsHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewsId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).DeleteNewsHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NewsService_DeleteNewsHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).DeleteNewsHandler(ctx, req.(*NewsId))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_ListNewsHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).ListNewsHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NewsService_ListNewsHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).ListNewsHandler(ctx, req.(*GetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NewsService_ServiceDesc is the grpc.ServiceDesc for NewsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NewsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "data.NewsService",
	HandlerType: (*NewsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNewsHandler",
			Handler:    _NewsService_CreateNewsHandler_Handler,
		},
		{
			MethodName: "ShowNewsHandler",
			Handler:    _NewsService_ShowNewsHandler_Handler,
		},
		{
			MethodName: "UpdateNewsHandler",
			Handler:    _NewsService_UpdateNewsHandler_Handler,
		},
		{
			MethodName: "DeleteNewsHandler",
			Handler:    _NewsService_DeleteNewsHandler_Handler,
		},
		{
			MethodName: "ListNewsHandler",
			Handler:    _NewsService_ListNewsHandler_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "news.proto",
}
