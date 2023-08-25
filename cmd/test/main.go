package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zcong1993/grpc-go-beyond/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type echoClient struct {
	client pb.HelloClient
}

func (ec *echoClient) Any(ctx context.Context, req *pb.EchoRequest) (*pb.EchoRequest, error) {
	resp, err := ec.client.Any(ctx, req)
	if err != nil {
		return nil, err
	}

	r := new(pb.EchoRequest)
	err = resp.Any.UnmarshalTo(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func main() {
	msg := &pb.EchoRequest{
		Message: "xxx",
	}
	ap, err := anypb.New(msg)
	check(err)
	ta := &pb.TestAny{
		Message:   "test",
		Any:       ap,
		TestOneof: &pb.TestAny_Name{Name: "test"},
	}

	oj, _ := json.Marshal(ta)
	pj, _ := protojson.Marshal(ta)
	fmt.Println(string(oj), string(pj))

	ta = &pb.TestAny{
		Message:   "test",
		Any:       ap,
		TestOneof: &pb.TestAny_EchoRequest{EchoRequest: msg},
	}

	oj, _ = json.Marshal(ta)
	pj, _ = protojson.Marshal(ta)
	fmt.Println(string(oj), string(pj))

	ta = &pb.TestAny{
		Message:   "test",
		Any:       ap,
		TestOneof: &pb.TestAny_OneofAny{OneofAny: ap},
	}

	oj, _ = json.Marshal(ta)
	pj, _ = protojson.Marshal(ta)
	fmt.Println(string(oj), string(pj))

	//bt, err := proto.Marshal(msg)
	//check(err)
	//var msg2 test.TestRequest
	//err = proto.Unmarshal(bt, &msg2)
	//check(err)
	//fmt.Println(msg2.String())
	//
	//ap, err := anypb.New(msg)
	//check(err)
	//fmt.Println(ap.String(), ap.TypeUrl)
	//var msg3 test.TestRequest
	//err = proto.Unmarshal(ap.Value, &msg3)
	//check(err)
	//fmt.Println(msg3.String())
	//
	//err = ap.UnmarshalTo(&msg3)
	//check(err)
	//fmt.Println(msg3.String())
}
