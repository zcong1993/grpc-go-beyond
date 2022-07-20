package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/zcong1993/grpc-go-beyond/internal/client"
	"github.com/zcong1993/grpc-go-beyond/pb"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", "localhost:8888", "server addr")
	method := flag.String("method", "Echo", "test method")
	clientType := flag.String("type", "stream", "client type: stream | generic | default")

	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	checkErr(err)

	var c pb.HelloClient

	switch *clientType {
	case "stream":
		c = client.NewStreamClient(conn)
	case "generic":
		c = client.NewGenericClient(conn)
	default:
		c = pb.NewHelloClient(conn)
	}

	m := map[string]func(c pb.HelloClient){
		"Echo":         testEcho,
		"ServerStream": testServerStream,
		"ClientStream": testClientStream,
		"DuplexStream": testDuplexStream,
	}

	h, ok := m[*method]
	if !ok {
		log.Fatal("invalid method")
	}

	h(c)
}

func testEcho(c pb.HelloClient) {
	ctx := context.Background()
	req := &pb.EchoRequest{Message: "test"}

	for {
		resp, err := c.Echo(ctx, req)
		checkErr(err)
		fmt.Println("send: ", req)
		fmt.Println("recv: ", resp)
		time.Sleep(time.Second * 5)
	}
}

func testServerStream(c pb.HelloClient) {
	ctx := context.Background()
	req := &pb.EchoRequest{Message: "test"}

	ss, err := c.ServerStream(ctx, req)
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

func testClientStream(c pb.HelloClient) {
	ctx := context.Background()

	s, err := c.ClientStream(ctx)
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

func testDuplexStream(c pb.HelloClient) {
	ctx := context.Background()

	s, err := c.DuplexStream(ctx)
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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
