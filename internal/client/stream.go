package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zcong1993/grpc-go-beyond/pb"
)

type StreamClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamClient(cc grpc.ClientConnInterface) *StreamClient {
	return &StreamClient{cc: cc}
}

func (s *StreamClient) Echo(ctx context.Context, in *pb.EchoRequest, opts ...grpc.CallOption) (*pb.EchoRequest, error) {
	cs, err := s.cc.NewStream(ctx, &grpc.StreamDesc{
		ServerStreams: false,
		ClientStreams: false,
	}, "/proto.Hello/Echo", opts...)

	if err != nil {
		return nil, err
	}

	err = cs.SendMsg(in)

	if err != nil {
		return nil, err
	}

	resp := new(pb.EchoRequest)
	if err := cs.RecvMsg(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type ss struct {
	grpc.ClientStream
}

func (s *ss) Recv() (*pb.EchoRequest, error) {
	resp := new(pb.EchoRequest)
	if err := s.ClientStream.RecvMsg(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *StreamClient) ServerStream(ctx context.Context, in *pb.EchoRequest, opts ...grpc.CallOption) (pb.Hello_ServerStreamClient, error) {
	cs, err := s.cc.NewStream(ctx, &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: false,
	}, "/proto.Hello/ServerStream", opts...)

	if err != nil {
		return nil, err
	}

	x := &ss{ClientStream: cs}

	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type css struct {
	grpc.ClientStream
}

func (x *css) Send(m *pb.EchoRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *css) CloseAndRecv() (*pb.EchoRequest, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(pb.EchoRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *StreamClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (pb.Hello_ClientStreamClient, error) {
	cs, err := s.cc.NewStream(ctx, &grpc.StreamDesc{
		ServerStreams: false,
		ClientStreams: true,
	}, "/proto.Hello/ClientStream", opts...)

	if err != nil {
		return nil, err
	}
	x := &css{cs}
	return x, nil
}

type dsc struct {
	grpc.ClientStream
}

func (x *dsc) Send(m *pb.EchoRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dsc) Recv() (*pb.EchoRequest, error) {
	m := new(pb.EchoRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *StreamClient) DuplexStream(ctx context.Context, opts ...grpc.CallOption) (pb.Hello_DuplexStreamClient, error) {
	cs, err := s.cc.NewStream(ctx, &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}, "/proto.Hello/DuplexStream", opts...)

	if err != nil {
		return nil, err
	}
	x := &dsc{cs}
	return x, nil
}
