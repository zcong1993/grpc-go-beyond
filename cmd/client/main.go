package main

import (
	"flag"
	"log"

	"github.com/zcong1993/grpc-go-beyond/internal/clienttest"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zcong1993/grpc-go-beyond/internal/client"
	"github.com/zcong1993/grpc-go-beyond/pb"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", "localhost:8888", "server addr")
	method := flag.String("method", "Echo", "test method: Echo | ServerStream | ClientStream | DuplexStream")
	clientType := flag.String("type", "stream", "client type: stream | raw | default")

	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	checkErr(err)

	var tester clienttest.Tester

	switch *clientType {
	case "stream":
		c := client.NewStreamClient(conn)
		tester = clienttest.NewHelloClientTester(c)
	case "raw":
		tester = clienttest.NewRawTester(conn)
	default:
		c := pb.NewHelloClient(conn)
		tester = clienttest.NewHelloClientTester(c)
	}

	switch *method {
	case "Echo":
		tester.TestEcho()
	case "ServerStream":
		tester.TestServerStream()
	case "ClientStream":
		tester.TestClientStream()
	case "DuplexStream":
		tester.TestDuplexStream()
	default:
		log.Fatal("invalid method")
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
