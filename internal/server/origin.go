package server

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zcong1993/grpc-go-beyond/pb"
)

var (
	header  = metadata.New(map[string]string{"x-test-header": "test"})
	trailer = metadata.New(map[string]string{"x-test-trailer": "test"})
)

type HelloServer struct {
	pb.UnimplementedHelloServer
}

func (h *HelloServer) Echo(ctx context.Context, request *pb.EchoRequest) (*pb.EchoRequest, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Println("metadata: ", md)
	}
	fmt.Println("recv: ", request)
	fmt.Println("send: ", request)
	err := grpc.SendHeader(ctx, header)
	if err != nil {
		return nil, err
	}
	err = grpc.SetTrailer(ctx, trailer)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (h *HelloServer) ServerStream(request *pb.EchoRequest, server pb.Hello_ServerStreamServer) error {
	if md, ok := metadata.FromIncomingContext(server.Context()); ok {
		fmt.Println("metadata: ", md)
	}
	fmt.Println("recv: ", request)

	err := server.SendHeader(header)
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		err := server.Send(request)
		if err != nil {
			return err
		}
		fmt.Println("send: ", request)
	}
	server.SetTrailer(trailer)

	return nil
}

func (h *HelloServer) ClientStream(server pb.Hello_ClientStreamServer) error {
	if md, ok := metadata.FromIncomingContext(server.Context()); ok {
		fmt.Println("metadata: ", md)
	}

	var last *pb.EchoRequest
	for {
		req, err := server.Recv()
		if err == io.EOF {
			fmt.Println("EOF recv")
			break
		} else if err != nil {
			return err
		}
		fmt.Println("recv: ", req)
		last = req
	}
	err := server.SendHeader(header)
	if err != nil {
		return err
	}
	fmt.Println("send: ", last)
	server.SetTrailer(trailer)
	return server.SendAndClose(last)
}

func (h *HelloServer) DuplexStream(server pb.Hello_DuplexStreamServer) error {
	if md, ok := metadata.FromIncomingContext(server.Context()); ok {
		fmt.Println("metadata: ", md)
	}

	for i := 0; ; i++ {
		req, err := server.Recv()
		if err == io.EOF {
			fmt.Println("EOF recv")
			break
		} else if err != nil {
			return err
		}

		fmt.Println("recv: ", req)
		if i == 0 {
			err := server.SendHeader(header)
			if err != nil {
				return err
			}
		}

		err2 := server.Send(req)
		if err2 != nil {
			return err2
		}
		fmt.Println("send: ", req)
	}

	server.SetTrailer(trailer)

	return nil
}
