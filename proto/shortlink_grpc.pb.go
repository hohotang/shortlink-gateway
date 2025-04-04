// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: proto/shortlink.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	URLService_ShortenURL_FullMethodName = "/shortlink.URLService/ShortenURL"
	URLService_ExpandURL_FullMethodName  = "/shortlink.URLService/ExpandURL"
)

// URLServiceClient is the client API for URLService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// URLService provides URL shortening and expansion functionality
type URLServiceClient interface {
	// ShortenURL creates a short URL from the original URL
	ShortenURL(ctx context.Context, in *ShortenURLRequest, opts ...grpc.CallOption) (*ShortenURLResponse, error)
	// ExpandURL resolves a short URL to its original URL
	ExpandURL(ctx context.Context, in *ExpandURLRequest, opts ...grpc.CallOption) (*ExpandURLResponse, error)
}

type uRLServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewURLServiceClient(cc grpc.ClientConnInterface) URLServiceClient {
	return &uRLServiceClient{cc}
}

func (c *uRLServiceClient) ShortenURL(ctx context.Context, in *ShortenURLRequest, opts ...grpc.CallOption) (*ShortenURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenURLResponse)
	err := c.cc.Invoke(ctx, URLService_ShortenURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLServiceClient) ExpandURL(ctx context.Context, in *ExpandURLRequest, opts ...grpc.CallOption) (*ExpandURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExpandURLResponse)
	err := c.cc.Invoke(ctx, URLService_ExpandURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLServiceServer is the server API for URLService service.
// All implementations must embed UnimplementedURLServiceServer
// for forward compatibility.
//
// URLService provides URL shortening and expansion functionality
type URLServiceServer interface {
	// ShortenURL creates a short URL from the original URL
	ShortenURL(context.Context, *ShortenURLRequest) (*ShortenURLResponse, error)
	// ExpandURL resolves a short URL to its original URL
	ExpandURL(context.Context, *ExpandURLRequest) (*ExpandURLResponse, error)
	mustEmbedUnimplementedURLServiceServer()
}

// UnimplementedURLServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedURLServiceServer struct{}

func (UnimplementedURLServiceServer) ShortenURL(context.Context, *ShortenURLRequest) (*ShortenURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenURL not implemented")
}
func (UnimplementedURLServiceServer) ExpandURL(context.Context, *ExpandURLRequest) (*ExpandURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExpandURL not implemented")
}
func (UnimplementedURLServiceServer) mustEmbedUnimplementedURLServiceServer() {}
func (UnimplementedURLServiceServer) testEmbeddedByValue()                    {}

// UnsafeURLServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLServiceServer will
// result in compilation errors.
type UnsafeURLServiceServer interface {
	mustEmbedUnimplementedURLServiceServer()
}

func RegisterURLServiceServer(s grpc.ServiceRegistrar, srv URLServiceServer) {
	// If the following call pancis, it indicates UnimplementedURLServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&URLService_ServiceDesc, srv)
}

func _URLService_ShortenURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServiceServer).ShortenURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLService_ShortenURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServiceServer).ShortenURL(ctx, req.(*ShortenURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLService_ExpandURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExpandURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServiceServer).ExpandURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLService_ExpandURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServiceServer).ExpandURL(ctx, req.(*ExpandURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLService_ServiceDesc is the grpc.ServiceDesc for URLService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortlink.URLService",
	HandlerType: (*URLServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenURL",
			Handler:    _URLService_ShortenURL_Handler,
		},
		{
			MethodName: "ExpandURL",
			Handler:    _URLService_ExpandURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortlink.proto",
}
