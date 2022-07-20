package server

import (
	"fmt"
	"io"

	"github.com/zcong1993/grpc-go-beyond/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Handler() grpc.StreamHandler {
	return func(srv interface{}, serverStream grpc.ServerStream) error {
		fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
		if !ok {
			return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
		}

		s := stream{serverStream: serverStream}

		switch fullMethodName {
		case "/proto.Hello/Echo":
			return s.handleEcho()
		case "/proto.Hello/ServerStream":
			return s.handleServerStream()
		case "/proto.Hello/ClientStream":
			return s.handleClientStream()
		case "/proto.Hello/DuplexStream":
			return s.handleDuplexStream()
		default:
			return status.Errorf(codes.Internal, "method not exists")
		}
	}
}

type stream struct {
	serverStream grpc.ServerStream
}

func (s *stream) handleEcho() error {
	var req pb.EchoRequest
	err := s.serverStream.RecvMsg(&req)
	if err != nil {
		return err
	}
	return s.serverStream.SendMsg(&req)
}

func (s *stream) handleServerStream() error {
	var req pb.EchoRequest
	err := s.serverStream.RecvMsg(&req)
	if err != nil {
		return err
	}
	fmt.Println("recv: ", &req)

	for i := 0; i < 5; i++ {
		err := s.serverStream.SendMsg(&req)
		if err != nil {
			return err
		}
		fmt.Println("send: ", &req)
	}

	return nil
}

func (s *stream) handleClientStream() error {
	var last *pb.EchoRequest
	for {
		req := new(pb.EchoRequest)
		err := s.serverStream.RecvMsg(req)
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
	return s.serverStream.SendMsg(last)
}

func (s *stream) handleDuplexStream() error {
	for {
		req := new(pb.EchoRequest)
		err := s.serverStream.RecvMsg(req)
		if err == io.EOF {
			fmt.Println("EOF recv")
			break
		} else if err != nil {
			return err
		}

		fmt.Println("recv: ", req)

		err2 := s.serverStream.SendMsg(req)
		if err2 != nil {
			return err2
		}
		fmt.Println("send: ", req)
	}

	return nil
}
