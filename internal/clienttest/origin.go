package clienttest

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/zcong1993/grpc-go-beyond/pb"
)

type HelloClientTester struct {
	c   pb.HelloClient
	ctx context.Context
}

func NewHelloClientTester(c pb.HelloClient) *HelloClientTester {
	return &HelloClientTester{
		c:   c,
		ctx: context.Background(),
	}
}

func (h *HelloClientTester) TestEcho() {
	req := &pb.EchoRequest{Message: "test"}

	for {
		resp, err := h.c.Echo(h.ctx, req)
		checkErr(err)
		fmt.Println("send: ", req)
		fmt.Println("recv: ", resp)
		time.Sleep(time.Second * 5)
	}
}

func (h *HelloClientTester) TestServerStream() {
	req := &pb.EchoRequest{Message: "test"}

	ss, err := h.c.ServerStream(h.ctx, req)
	checkErr(err)
	fmt.Println("send: ", req)

	for {
		rr, err := ss.Recv()
		if err == io.EOF {
			fmt.Println("EOF recv")
			break
		}
		checkErr(err)

		fmt.Println("recv: ", rr)
	}
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
	fmt.Println("recv: ", resp)
}

func (h *HelloClientTester) TestDuplexStream() {
	s, err := h.c.DuplexStream(h.ctx)
	checkErr(err)

	req := &pb.EchoRequest{Message: "test"}

	ch := make(chan struct{})

	go func() {
		for {
			rr, err := s.Recv()
			if err == io.EOF {
				fmt.Println("EOF recv")
				break
			}
			checkErr(err)

			fmt.Println("recv: ", rr)
		}
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
