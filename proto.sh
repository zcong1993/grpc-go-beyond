#!/bin/bash -x
#set -e

OUT_DIR=./pb

rm -rf $OUT_DIR/*.go

PROTO_PREFIX=github.com/zcong1993/grpc-go-beyond

protoc --go_out=. --go_opt=module=$PROTO_PREFIX \
    --go-grpc_out=. --go-grpc_opt=module=$PROTO_PREFIX \
    proto/*.proto
