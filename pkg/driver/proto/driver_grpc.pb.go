// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: pkg/driver/proto/driver.proto

package proto

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

// RavelDriverClient is the client API for RavelDriver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RavelDriverClient interface {
	StartVM(ctx context.Context, in *StartVMRequest, opts ...grpc.CallOption) (*StartVMResponse, error)
	StopVM(ctx context.Context, in *StopVMRequest, opts ...grpc.CallOption) (*Empty, error)
}

type ravelDriverClient struct {
	cc grpc.ClientConnInterface
}

func NewRavelDriverClient(cc grpc.ClientConnInterface) RavelDriverClient {
	return &ravelDriverClient{cc}
}

func (c *ravelDriverClient) StartVM(ctx context.Context, in *StartVMRequest, opts ...grpc.CallOption) (*StartVMResponse, error) {
	out := new(StartVMResponse)
	err := c.cc.Invoke(ctx, "/proto.RavelDriver/StartVM", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ravelDriverClient) StopVM(ctx context.Context, in *StopVMRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.RavelDriver/StopVM", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RavelDriverServer is the server API for RavelDriver service.
// All implementations must embed UnimplementedRavelDriverServer
// for forward compatibility
type RavelDriverServer interface {
	StartVM(context.Context, *StartVMRequest) (*StartVMResponse, error)
	StopVM(context.Context, *StopVMRequest) (*Empty, error)
	mustEmbedUnimplementedRavelDriverServer()
}

// UnimplementedRavelDriverServer must be embedded to have forward compatible implementations.
type UnimplementedRavelDriverServer struct {
}

func (UnimplementedRavelDriverServer) StartVM(context.Context, *StartVMRequest) (*StartVMResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartVM not implemented")
}
func (UnimplementedRavelDriverServer) StopVM(context.Context, *StopVMRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopVM not implemented")
}
func (UnimplementedRavelDriverServer) mustEmbedUnimplementedRavelDriverServer() {}

// UnsafeRavelDriverServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RavelDriverServer will
// result in compilation errors.
type UnsafeRavelDriverServer interface {
	mustEmbedUnimplementedRavelDriverServer()
}

func RegisterRavelDriverServer(s grpc.ServiceRegistrar, srv RavelDriverServer) {
	s.RegisterService(&RavelDriver_ServiceDesc, srv)
}

func _RavelDriver_StartVM_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartVMRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RavelDriverServer).StartVM(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RavelDriver/StartVM",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RavelDriverServer).StartVM(ctx, req.(*StartVMRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RavelDriver_StopVM_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopVMRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RavelDriverServer).StopVM(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RavelDriver/StopVM",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RavelDriverServer).StopVM(ctx, req.(*StopVMRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RavelDriver_ServiceDesc is the grpc.ServiceDesc for RavelDriver service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RavelDriver_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.RavelDriver",
	HandlerType: (*RavelDriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartVM",
			Handler:    _RavelDriver_StartVM_Handler,
		},
		{
			MethodName: "StopVM",
			Handler:    _RavelDriver_StopVM_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/driver/proto/driver.proto",
}
