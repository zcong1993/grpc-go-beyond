package server

import (
	"context"
	"fmt"
	"io"

	"github.com/zcong1993/grpc-go-beyond/pb"
)

type HelloServer struct {
	pb.UnimplementedHelloServer
}

func (h *HelloServer) Echo(ctx context.Context, request *pb.EchoRequest) (*pb.EchoRequest, error) {
	fmt.Println("recv: ", request)
	fmt.Println("send: ", request)
	return request, nil
}

func (h *HelloServer) ServerStream(request *pb.EchoRequest, server pb.Hello_ServerStreamServer) error {
	fmt.Println("recv: ", request)

	for i := 0; i < 5; i++ {
		err := server.Send(request)
		if err != nil {
			return err
		}
		fmt.Println("send: ", request)
	}

	return nil
}

func (h *HelloServer) ClientStream(server pb.Hello_ClientStreamServer) error {
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
	fmt.Println("send: ", last)
	return server.SendAndClose(last)
}

func (h *HelloServer) DuplexStream(server pb.Hello_DuplexStreamServer) error {
	for {
		req, err := server.Recv()
		if err == io.EOF {
			fmt.Println("EOF recv")
			break
		} else if err != nil {
			return err
		}

		fmt.Println("recv: ", req)

		err2 := server.Send(req)
		if err2 != nil {
			return err2
		}
		fmt.Println("send: ", req)
	}

	return nil
}
