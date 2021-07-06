// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package test

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

// MessagesServiceClient is the client API for MessagesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessagesServiceClient interface {
	// Sends a greeting
	Chat(ctx context.Context, opts ...grpc.CallOption) (MessagesService_ChatClient, error)
}

type messagesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessagesServiceClient(cc grpc.ClientConnInterface) MessagesServiceClient {
	return &messagesServiceClient{cc}
}

func (c *messagesServiceClient) Chat(ctx context.Context, opts ...grpc.CallOption) (MessagesService_ChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessagesService_ServiceDesc.Streams[0], "/helloworld.MessagesService/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &messagesServiceChatClient{stream}
	return x, nil
}

type MessagesService_ChatClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type messagesServiceChatClient struct {
	grpc.ClientStream
}

func (x *messagesServiceChatClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messagesServiceChatClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessagesServiceServer is the server API for MessagesService service.
// All implementations must embed UnimplementedMessagesServiceServer
// for forward compatibility
type MessagesServiceServer interface {
	// Sends a greeting
	Chat(MessagesService_ChatServer) error
	mustEmbedUnimplementedMessagesServiceServer()
}

// UnimplementedMessagesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessagesServiceServer struct {
}

func (UnimplementedMessagesServiceServer) Chat(MessagesService_ChatServer) error {
	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
}
func (UnimplementedMessagesServiceServer) mustEmbedUnimplementedMessagesServiceServer() {}

// UnsafeMessagesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessagesServiceServer will
// result in compilation errors.
type UnsafeMessagesServiceServer interface {
	mustEmbedUnimplementedMessagesServiceServer()
}

func RegisterMessagesServiceServer(s grpc.ServiceRegistrar, srv MessagesServiceServer) {
	s.RegisterService(&MessagesService_ServiceDesc, srv)
}

func _MessagesService_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessagesServiceServer).Chat(&messagesServiceChatServer{stream})
}

type MessagesService_ChatServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type messagesServiceChatServer struct {
	grpc.ServerStream
}

func (x *messagesServiceChatServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messagesServiceChatServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessagesService_ServiceDesc is the grpc.ServiceDesc for MessagesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessagesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.MessagesService",
	HandlerType: (*MessagesServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _MessagesService_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/greeter.proto",
}
