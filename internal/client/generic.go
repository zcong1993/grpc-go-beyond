package client

import (
	"context"
	"fmt"
	"io"

	"github.com/zcong1993/grpc-go-beyond/pb"
	"google.golang.org/grpc"
)

var desc = &grpc.StreamDesc{
	ServerStreams: true,
	ClientStreams: true,
}

type GenericClient struct {
	cc grpc.ClientConnInterface
}

func NewGenericClient(cc grpc.ClientConnInterface) *GenericClient {
	return &GenericClient{cc: cc}
}

func (g *GenericClient) Echo(ctx context.Context, in *pb.EchoRequest, opts ...grpc.CallOption) (*pb.EchoRequest, error) {
	cs, err := g.cc.NewStream(ctx, desc, "/proto.Hello/Echo", opts...)

	if err != nil {
		return nil, err
	}

	err = cs.SendMsg(in)
	if err != nil {
		return nil, err
	}
	fmt.Println("generic send: ", in)

	err = cs.CloseSend()

	if err != nil {
		return nil, err
	}

	resp := new(pb.EchoRequest)

	for {
		err = cs.RecvMsg(resp)
		if err == io.EOF {
			fmt.Println("generic EOF recv")
			break
		} else if err != nil {
			return nil, err
		}
		fmt.Println("generic recv: ", resp)
	}

	return resp, nil
}

func (g *GenericClient) ServerStream(ctx context.Context, in *pb.EchoRequest, opts ...grpc.CallOption) (pb.Hello_ServerStreamClient, error) {
	cs, err := g.cc.NewStream(ctx, desc, "/proto.Hello/ServerStream", opts...)
	if err != nil {
		return nil, err
	}

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

func (g *GenericClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (pb.Hello_ClientStreamClient, error) {
	cs, err := g.cc.NewStream(ctx, desc, "/proto.Hello/ClientStream", opts...)
	if err != nil {
		return nil, err
	}

	x := &css{cs}
	return x, nil
}

func (g *GenericClient) DuplexStream(ctx context.Context, opts ...grpc.CallOption) (pb.Hello_DuplexStreamClient, error) {
	cs, err := g.cc.NewStream(ctx, desc, "/proto.Hello/DuplexStream", opts...)
	if err != nil {
		return nil, err
	}

	x := &dsc{cs}
	return x, nil
}
