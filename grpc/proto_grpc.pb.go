// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: grpc/proto.proto

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

const (
	StreamingService_GetChatMessageStreaming_FullMethodName = "/simpleGuide.StreamingService/GetChatMessageStreaming"
	StreamingService_SendChatMessage_FullMethodName         = "/simpleGuide.StreamingService/SendChatMessage"
)

// StreamingServiceClient is the client API for StreamingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamingServiceClient interface {
	// rpc Join(JoinRequest) returns (JoinResponse) {}
	GetChatMessageStreaming(ctx context.Context, in *Connect, opts ...grpc.CallOption) (StreamingService_GetChatMessageStreamingClient, error)
	SendChatMessage(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*Empty, error)
}

type streamingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamingServiceClient(cc grpc.ClientConnInterface) StreamingServiceClient {
	return &streamingServiceClient{cc}
}

func (c *streamingServiceClient) GetChatMessageStreaming(ctx context.Context, in *Connect, opts ...grpc.CallOption) (StreamingService_GetChatMessageStreamingClient, error) {
	stream, err := c.cc.NewStream(ctx, &StreamingService_ServiceDesc.Streams[0], StreamingService_GetChatMessageStreaming_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &streamingServiceGetChatMessageStreamingClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamingService_GetChatMessageStreamingClient interface {
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type streamingServiceGetChatMessageStreamingClient struct {
	grpc.ClientStream
}

func (x *streamingServiceGetChatMessageStreamingClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamingServiceClient) SendChatMessage(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, StreamingService_SendChatMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StreamingServiceServer is the server API for StreamingService service.
// All implementations must embed UnimplementedStreamingServiceServer
// for forward compatibility
type StreamingServiceServer interface {
	// rpc Join(JoinRequest) returns (JoinResponse) {}
	GetChatMessageStreaming(*Connect, StreamingService_GetChatMessageStreamingServer) error
	SendChatMessage(context.Context, *ChatMessage) (*Empty, error)
	mustEmbedUnimplementedStreamingServiceServer()
}

// UnimplementedStreamingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStreamingServiceServer struct {
}

func (UnimplementedStreamingServiceServer) GetChatMessageStreaming(*Connect, StreamingService_GetChatMessageStreamingServer) error {
	return status.Errorf(codes.Unimplemented, "method GetChatMessageStreaming not implemented")
}
func (UnimplementedStreamingServiceServer) SendChatMessage(context.Context, *ChatMessage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendChatMessage not implemented")
}
func (UnimplementedStreamingServiceServer) mustEmbedUnimplementedStreamingServiceServer() {}

// UnsafeStreamingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamingServiceServer will
// result in compilation errors.
type UnsafeStreamingServiceServer interface {
	mustEmbedUnimplementedStreamingServiceServer()
}

func RegisterStreamingServiceServer(s grpc.ServiceRegistrar, srv StreamingServiceServer) {
	s.RegisterService(&StreamingService_ServiceDesc, srv)
}

func _StreamingService_GetChatMessageStreaming_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Connect)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamingServiceServer).GetChatMessageStreaming(m, &streamingServiceGetChatMessageStreamingServer{stream})
}

type StreamingService_GetChatMessageStreamingServer interface {
	Send(*ChatMessage) error
	grpc.ServerStream
}

type streamingServiceGetChatMessageStreamingServer struct {
	grpc.ServerStream
}

func (x *streamingServiceGetChatMessageStreamingServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _StreamingService_SendChatMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamingServiceServer).SendChatMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StreamingService_SendChatMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamingServiceServer).SendChatMessage(ctx, req.(*ChatMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// StreamingService_ServiceDesc is the grpc.ServiceDesc for StreamingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "simpleGuide.StreamingService",
	HandlerType: (*StreamingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendChatMessage",
			Handler:    _StreamingService_SendChatMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetChatMessageStreaming",
			Handler:       _StreamingService_GetChatMessageStreaming_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc/proto.proto",
}
