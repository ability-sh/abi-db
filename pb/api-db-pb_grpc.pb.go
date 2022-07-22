// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: pb/api-db-pb.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	Get(ctx context.Context, in *GetTask, opts ...grpc.CallOption) (*GetResult, error)
	Put(ctx context.Context, in *PutTask, opts ...grpc.CallOption) (*PutResult, error)
	Merge(ctx context.Context, in *MergeTask, opts ...grpc.CallOption) (*MergeResult, error)
	Del(ctx context.Context, in *DelTask, opts ...grpc.CallOption) (*DelResult, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Get(ctx context.Context, in *GetTask, opts ...grpc.CallOption) (*GetResult, error) {
	out := new(GetResult)
	err := c.cc.Invoke(ctx, "/db.Service/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Put(ctx context.Context, in *PutTask, opts ...grpc.CallOption) (*PutResult, error) {
	out := new(PutResult)
	err := c.cc.Invoke(ctx, "/db.Service/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Merge(ctx context.Context, in *MergeTask, opts ...grpc.CallOption) (*MergeResult, error) {
	out := new(MergeResult)
	err := c.cc.Invoke(ctx, "/db.Service/Merge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Del(ctx context.Context, in *DelTask, opts ...grpc.CallOption) (*DelResult, error) {
	out := new(DelResult)
	err := c.cc.Invoke(ctx, "/db.Service/Del", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations should embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	Get(context.Context, *GetTask) (*GetResult, error)
	Put(context.Context, *PutTask) (*PutResult, error)
	Merge(context.Context, *MergeTask) (*MergeResult, error)
	Del(context.Context, *DelTask) (*DelResult, error)
}

// UnimplementedServiceServer should be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) Get(context.Context, *GetTask) (*GetResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedServiceServer) Put(context.Context, *PutTask) (*PutResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedServiceServer) Merge(context.Context, *MergeTask) (*MergeResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Merge not implemented")
}
func (UnimplementedServiceServer) Del(context.Context, *DelTask) (*DelResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/db.Service/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Get(ctx, req.(*GetTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/db.Service/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Put(ctx, req.(*PutTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Merge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Merge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/db.Service/Merge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Merge(ctx, req.(*MergeTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/db.Service/Del",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Del(ctx, req.(*DelTask))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "db.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Service_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _Service_Put_Handler,
		},
		{
			MethodName: "Merge",
			Handler:    _Service_Merge_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _Service_Del_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/api-db-pb.proto",
}
