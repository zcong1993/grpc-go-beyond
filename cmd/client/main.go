package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/zcong1993/grpc-go-beyond/internal/sign"
	"github.com/zcong1993/grpc-go-beyond/pb/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/zcong1993/grpc-go-beyond/pb"
)

func signMw(app, key string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		r, ok := req.(proto.Message)
		if !ok {
			return fmt.Errorf("req must be proto.Message")
		}
		rawReq, err := proto.Marshal(r)
		if err != nil {
			return err
		}
		rawReq = append(rawReq, []byte(app)...)
		s := sign.Sign(key, rawReq)
		fmt.Println("b64: ", sign.B64(rawReq))
		fmt.Println("sign: ", s)
		md := metadata.Pairs("app", app, "sign", s)
		//md := metadata.Pairs()
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func main() {
	addr := flag.String("addr", "localhost:9999", "server addr")
	app := flag.String("app", "biz1", "app name")
	method := flag.String("method", "Echo", "method name")
	signKey := flag.String("sign-key", "biz1", "sign key")

	flag.Parse()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(signMw(*app, *signKey))}

	conn, err := grpc.Dial(*addr, opts...)
	checkErr(err)

	c := pb.NewHelloClient(conn)
	tc := test.NewTestClient(conn)

	ctx := context.Background()
	fmt.Println("call ", *method)

	switch *method {
	case "Echo":
		resp, err := c.Echo(ctx, &pb.EchoRequest{Message: "hello"})
		checkErr(err)
		fmt.Println("resp: ", resp)
	case "Echo2":
		resp, err := c.Echo2(ctx, &pb.EchoRequest{Message: "hello"})
		checkErr(err)
		fmt.Println("resp: ", resp)
	case "Test":
		resp, err := tc.Test(ctx, &test.TestRequest{Name: "hello"})
		checkErr(err)
		fmt.Println("resp: ", resp)
	case "Test1":
		resp, err := tc.Test1(ctx, &test.TestRequest{Name: "hello"})
		checkErr(err)
		fmt.Println("resp: ", resp)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
