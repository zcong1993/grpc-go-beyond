package server

import (
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/zcong1993/grpc-go-beyond/pb"
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
	if md, ok := metadata.FromIncomingContext(s.serverStream.Context()); ok {
		fmt.Println("metadata: ", md)
	}

	var req pb.EchoRequest
	err := s.serverStream.RecvMsg(&req)
	if err != nil {
		return err
	}
	err = s.serverStream.SendHeader(header)
	if err != nil {
		return err
	}

	defer s.serverStream.SetTrailer(trailer)

	return s.serverStream.SendMsg(&req)
}

func (s *stream) handleServerStream() error {
	if md, ok := metadata.FromIncomingContext(s.serverStream.Context()); ok {
		fmt.Println("metadata: ", md)
	}

	var req pb.EchoRequest
	err := s.serverStream.RecvMsg(&req)
	if err != nil {
		return err
	}
	fmt.Println("recv: ", &req)
	err = s.serverStream.SendHeader(header)
	if err != nil {
		return err
	}

	defer s.serverStream.SetTrailer(trailer)

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
	if md, ok := metadata.FromIncomingContext(s.serverStream.Context()); ok {
		fmt.Println("metadata: ", md)
	}

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

	err := s.serverStream.SendHeader(header)
	if err != nil {
		return err
	}
	fmt.Println("send: ", last)
	defer s.serverStream.SetTrailer(trailer)
	return s.serverStream.SendMsg(last)
}

func (s *stream) handleDuplexStream() error {
	if md, ok := metadata.FromIncomingContext(s.serverStream.Context()); ok {
		fmt.Println("metadata: ", md)
	}

	headerSent := false
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

		if !headerSent {
			headerSent = true
			err := s.serverStream.SendHeader(header)
			if err != nil {
				return err
			}
		}

		err2 := s.serverStream.SendMsg(req)
		if err2 != nil {
			return err2
		}
		fmt.Println("send: ", req)
	}

	s.serverStream.SetTrailer(trailer)

	return nil
}
