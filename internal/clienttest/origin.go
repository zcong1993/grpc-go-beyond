package clienttest

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zcong1993/grpc-go-beyond/pb"
)

type HelloClientTester struct {
	c   pb.HelloClient
	ctx context.Context
}

func NewHelloClientTester(c pb.HelloClient) *HelloClientTester {
	return &HelloClientTester{
		c:   c,
		ctx: metadata.AppendToOutgoingContext(context.Background(), "aaa", "bbb"),
	}
}

func (h *HelloClientTester) TestEcho() {
	req := &pb.EchoRequest{Message: "test"}

	var header, trailer metadata.MD
	resp, err := h.c.Echo(h.ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
	checkErr(err)
	fmt.Println("send: ", req)
	fmt.Println("recv: ", resp)
	fmt.Println("header: ", header)
	fmt.Println("trailer: ", trailer)
}

func (h *HelloClientTester) TestServerStream() {
	req := &pb.EchoRequest{Message: "test"}

	ss, err := h.c.ServerStream(h.ctx, req)
	checkErr(err)
	fmt.Println("send: ", req)

	md, err := ss.Header()
	checkErr(err)
	fmt.Println("header: ", md)

	for {
		rr, err := ss.Recv()
		if err == io.EOF {
			fmt.Println("EOF recv")
			break
		}
		checkErr(err)

		fmt.Println("recv: ", rr)
	}

	fmt.Println("trailer: ", ss.Trailer())
}

func (h *HelloClientTester) TestClientStream() {
	s, err := h.c.ClientStream(h.ctx)
	checkErr(err)

	req := &pb.EchoRequest{Message: "test"}

	for i := 0; i < 5; i++ {
		err := s.Send(req)
		checkErr(err)
		fmt.Println("send: ", req)
	}

	resp, err := s.CloseAndRecv()
	checkErr(err)

	md, err := s.Header()
	checkErr(err)
	fmt.Println("header: ", md)

	fmt.Println("recv: ", resp)
	fmt.Println("trailer: ", s.Trailer())
}

func (h *HelloClientTester) TestDuplexStream() {
	s, err := h.c.DuplexStream(h.ctx)
	checkErr(err)

	req := &pb.EchoRequest{Message: "test"}

	ch := make(chan struct{})

	go func() {
		md, err := s.Header()
		checkErr(err)
		fmt.Println("header: ", md)

		for {
			rr, err := s.Recv()
			if err == io.EOF {
				fmt.Println("EOF recv")
				break
			}
			checkErr(err)

			fmt.Println("recv: ", rr)
		}

		fmt.Println("trailer: ", s.Trailer())
		ch <- struct{}{}
	}()

	for i := 0; i < 5; i++ {
		err := s.Send(req)
		checkErr(err)
		fmt.Println("send: ", req)
		time.Sleep(time.Millisecond * 200)
	}
	err = s.CloseSend()
	checkErr(err)
	<-ch
}
