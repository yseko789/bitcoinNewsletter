// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: proto/summary.proto

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

// SummaryServiceClient is the client API for SummaryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SummaryServiceClient interface {
	ReadSummary(ctx context.Context, in *Date, opts ...grpc.CallOption) (*Summary, error)
}

type summaryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSummaryServiceClient(cc grpc.ClientConnInterface) SummaryServiceClient {
	return &summaryServiceClient{cc}
}

func (c *summaryServiceClient) ReadSummary(ctx context.Context, in *Date, opts ...grpc.CallOption) (*Summary, error) {
	out := new(Summary)
	err := c.cc.Invoke(ctx, "/bitcoinNewsletter.SummaryService/ReadSummary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SummaryServiceServer is the server API for SummaryService service.
// All implementations must embed UnimplementedSummaryServiceServer
// for forward compatibility
type SummaryServiceServer interface {
	ReadSummary(context.Context, *Date) (*Summary, error)
	mustEmbedUnimplementedSummaryServiceServer()
}

// UnimplementedSummaryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSummaryServiceServer struct {
}

func (UnimplementedSummaryServiceServer) ReadSummary(context.Context, *Date) (*Summary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadSummary not implemented")
}
func (UnimplementedSummaryServiceServer) mustEmbedUnimplementedSummaryServiceServer() {}

// UnsafeSummaryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SummaryServiceServer will
// result in compilation errors.
type UnsafeSummaryServiceServer interface {
	mustEmbedUnimplementedSummaryServiceServer()
}

func RegisterSummaryServiceServer(s grpc.ServiceRegistrar, srv SummaryServiceServer) {
	s.RegisterService(&SummaryService_ServiceDesc, srv)
}

func _SummaryService_ReadSummary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Date)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SummaryServiceServer).ReadSummary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bitcoinNewsletter.SummaryService/ReadSummary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SummaryServiceServer).ReadSummary(ctx, req.(*Date))
	}
	return interceptor(ctx, in, info, handler)
}

// SummaryService_ServiceDesc is the grpc.ServiceDesc for SummaryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SummaryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bitcoinNewsletter.SummaryService",
	HandlerType: (*SummaryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReadSummary",
			Handler:    _SummaryService_ReadSummary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/summary.proto",
}
