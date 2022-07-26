package clienttest

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zcong1993/grpc-go-beyond/pb"
)

var desc = &grpc.StreamDesc{
	ServerStreams: true,
	ClientStreams: true,
}

type RawTester struct {
	conn grpc.ClientConnInterface
	ctx  context.Context
}

func NewRawTester(conn grpc.ClientConnInterface) *RawTester {
	return &RawTester{
		conn: conn,
		ctx:  metadata.AppendToOutgoingContext(context.Background(), "aaa", "bbb"),
	}
}

func (r *RawTester) TestEcho() {
	cs, err := r.conn.NewStream(r.ctx, desc, "/proto.Hello/Echo")
	checkErr(err)

	req := &pb.EchoRequest{Message: "test"}

	err = cs.SendMsg(req)
	checkErr(err)
	fmt.Println("send: ", req)
	err = cs.CloseSend()
	checkErr(err)

	resp := new(pb.EchoRequest)

	md, err := cs.Header()
	checkErr(err)
	fmt.Println("header: ", md)

	err = cs.RecvMsg(resp)
	checkErr(err)
	fmt.Println("recv: ", resp)
	fmt.Println("trailer: ", cs.Trailer())
}

func (r *RawTester) TestServerStream() {
	cs, err := r.conn.NewStream(r.ctx, desc, "/proto.Hello/ServerStream")
	checkErr(err)

	req := &pb.EchoRequest{Message: "test"}

	err = cs.SendMsg(req)
	checkErr(err)
	fmt.Println("send: ", req)
	err = cs.CloseSend()
	checkErr(err)

	md, err := cs.Header()
	checkErr(err)
	fmt.Println("header: ", md)

	for {
		resp := new(pb.EchoRequest)
		err = cs.RecvMsg(resp)
		if err == io.EOF {
			fmt.Println("EOF recv")
			break
		}
		checkErr(err)
		fmt.Println("recv: ", resp)
	}

	fmt.Println("trailer: ", cs.Trailer())
}

func (r *RawTester) TestClientStream() {
	cs, err := r.conn.NewStream(r.ctx, desc, "/proto.Hello/ClientStream")
	checkErr(err)

	req := &pb.EchoRequest{Message: "test"}
	for i := 0; i < 5; i++ {
		err := cs.SendMsg(req)
		checkErr(err)
		fmt.Println("send: ", req)
	}

	err = cs.CloseSend()
	checkErr(err)

	md, err := cs.Header()
	checkErr(err)
	fmt.Println("header: ", md)

	resp := new(pb.EchoRequest)
	err = cs.RecvMsg(resp)
	checkErr(err)
	fmt.Println("recv: ", resp)
	fmt.Println("trailer: ", cs.Trailer())
}

func (r *RawTester) TestDuplexStream() {
	cs, err := r.conn.NewStream(r.ctx, desc, "/proto.Hello/DuplexStream")
	checkErr(err)

	req := &pb.EchoRequest{Message: "test"}

	ch := make(chan struct{})

	go func() {
		md, err := cs.Header()
		checkErr(err)
		fmt.Println("header: ", md)

		for {
			rr := new(pb.EchoRequest)
			err := cs.RecvMsg(rr)
			if err == io.EOF {
				fmt.Println("EOF recv")
				break
			}
			checkErr(err)

			fmt.Println("recv: ", rr)
		}

		fmt.Println("trailer: ", cs.Trailer())
		ch <- struct{}{}
	}()

	for i := 0; i < 5; i++ {
		err := cs.SendMsg(req)
		checkErr(err)
		fmt.Println("send: ", req)
		time.Sleep(time.Millisecond * 200)
	}
	err = cs.CloseSend()
	checkErr(err)
	<-ch
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
