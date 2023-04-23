package main

import (
	"fmt"
	"strings"

	"github.com/zcong1993/grpc-go-beyond/pb"
)

func main() {
	for i := 0; i < 10000; i++ {
		a := &pb.TestRequest{
			Message: "test",
			Name:    "xxx",
		} // 下面不加空行 space count=2, 加一行空行 space count=1

		fmt.Printf("%d, message: %s, space count:%d \n", i, a.String(), strings.Count(a.String(), " "))
		// space count=1 sample: 9999, message: message:"test" name:"xxx", space count:1
		// space count=2 sample: 9999, message: message:"test"  name:"xxx", space count:2
	}
}
