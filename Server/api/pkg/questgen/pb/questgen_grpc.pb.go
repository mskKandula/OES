// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: questgen.proto

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

// QuestGenServiceClient is the client API for QuestGenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuestGenServiceClient interface {
	// Unary
	QuestGen(ctx context.Context, in *QuestGenRequest, opts ...grpc.CallOption) (*QuestGenResponse, error)
}

type questGenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQuestGenServiceClient(cc grpc.ClientConnInterface) QuestGenServiceClient {
	return &questGenServiceClient{cc}
}

func (c *questGenServiceClient) QuestGen(ctx context.Context, in *QuestGenRequest, opts ...grpc.CallOption) (*QuestGenResponse, error) {
	out := new(QuestGenResponse)
	err := c.cc.Invoke(ctx, "/questgen.QuestGenService/QuestGen", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuestGenServiceServer is the server API for QuestGenService service.
// All implementations must embed UnimplementedQuestGenServiceServer
// for forward compatibility
type QuestGenServiceServer interface {
	// Unary
	QuestGen(context.Context, *QuestGenRequest) (*QuestGenResponse, error)
	mustEmbedUnimplementedQuestGenServiceServer()
}

// UnimplementedQuestGenServiceServer must be embedded to have forward compatible implementations.
type UnimplementedQuestGenServiceServer struct {
}

func (UnimplementedQuestGenServiceServer) QuestGen(context.Context, *QuestGenRequest) (*QuestGenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuestGen not implemented")
}
func (UnimplementedQuestGenServiceServer) mustEmbedUnimplementedQuestGenServiceServer() {}

// UnsafeQuestGenServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QuestGenServiceServer will
// result in compilation errors.
type UnsafeQuestGenServiceServer interface {
	mustEmbedUnimplementedQuestGenServiceServer()
}

func RegisterQuestGenServiceServer(s grpc.ServiceRegistrar, srv QuestGenServiceServer) {
	s.RegisterService(&QuestGenService_ServiceDesc, srv)
}

func _QuestGenService_QuestGen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestGenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestGenServiceServer).QuestGen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/questgen.QuestGenService/QuestGen",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestGenServiceServer).QuestGen(ctx, req.(*QuestGenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// QuestGenService_ServiceDesc is the grpc.ServiceDesc for QuestGenService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QuestGenService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "questgen.QuestGenService",
	HandlerType: (*QuestGenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QuestGen",
			Handler:    _QuestGenService_QuestGen_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "questgen.proto",
}
