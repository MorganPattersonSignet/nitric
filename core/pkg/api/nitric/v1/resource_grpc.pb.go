// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: proto/resource/v1/resource.proto

package v1

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

// ResourceServiceClient is the client API for ResourceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResourceServiceClient interface {
	// Declare a resource for the nitric application
	// At Deploy time this will create resources as part of the nitric stacks dependency graph
	// At runtime
	Declare(ctx context.Context, in *ResourceDeclareRequest, opts ...grpc.CallOption) (*ResourceDeclareResponse, error)
	// Retrieve details about a resource at runtime
	Details(ctx context.Context, in *ResourceDetailsRequest, opts ...grpc.CallOption) (*ResourceDetailsResponse, error)
}

type resourceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewResourceServiceClient(cc grpc.ClientConnInterface) ResourceServiceClient {
	return &resourceServiceClient{cc}
}

func (c *resourceServiceClient) Declare(ctx context.Context, in *ResourceDeclareRequest, opts ...grpc.CallOption) (*ResourceDeclareResponse, error) {
	out := new(ResourceDeclareResponse)
	err := c.cc.Invoke(ctx, "/nitric.resource.v1.ResourceService/Declare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceServiceClient) Details(ctx context.Context, in *ResourceDetailsRequest, opts ...grpc.CallOption) (*ResourceDetailsResponse, error) {
	out := new(ResourceDetailsResponse)
	err := c.cc.Invoke(ctx, "/nitric.resource.v1.ResourceService/Details", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceServiceServer is the server API for ResourceService service.
// All implementations must embed UnimplementedResourceServiceServer
// for forward compatibility
type ResourceServiceServer interface {
	// Declare a resource for the nitric application
	// At Deploy time this will create resources as part of the nitric stacks dependency graph
	// At runtime
	Declare(context.Context, *ResourceDeclareRequest) (*ResourceDeclareResponse, error)
	// Retrieve details about a resource at runtime
	Details(context.Context, *ResourceDetailsRequest) (*ResourceDetailsResponse, error)
	mustEmbedUnimplementedResourceServiceServer()
}

// UnimplementedResourceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedResourceServiceServer struct {
}

func (UnimplementedResourceServiceServer) Declare(context.Context, *ResourceDeclareRequest) (*ResourceDeclareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Declare not implemented")
}
func (UnimplementedResourceServiceServer) Details(context.Context, *ResourceDetailsRequest) (*ResourceDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Details not implemented")
}
func (UnimplementedResourceServiceServer) mustEmbedUnimplementedResourceServiceServer() {}

// UnsafeResourceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResourceServiceServer will
// result in compilation errors.
type UnsafeResourceServiceServer interface {
	mustEmbedUnimplementedResourceServiceServer()
}

func RegisterResourceServiceServer(s grpc.ServiceRegistrar, srv ResourceServiceServer) {
	s.RegisterService(&ResourceService_ServiceDesc, srv)
}

func _ResourceService_Declare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceDeclareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServiceServer).Declare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.resource.v1.ResourceService/Declare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServiceServer).Declare(ctx, req.(*ResourceDeclareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ResourceService_Details_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceDetailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServiceServer).Details(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.resource.v1.ResourceService/Details",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServiceServer).Details(ctx, req.(*ResourceDetailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ResourceService_ServiceDesc is the grpc.ServiceDesc for ResourceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ResourceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nitric.resource.v1.ResourceService",
	HandlerType: (*ResourceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Declare",
			Handler:    _ResourceService_Declare_Handler,
		},
		{
			MethodName: "Details",
			Handler:    _ResourceService_Details_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/resource/v1/resource.proto",
}
