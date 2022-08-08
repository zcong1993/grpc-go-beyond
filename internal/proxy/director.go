// Based on https://github.com/trusch/grpc-proxy
// Copyright Michal Witkowski. Licensed under Apache2 license: https://github.com/trusch/grpc-proxy/blob/master/LICENSE.txt

package proxy

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zcong1993/grpc-go-beyond/internal/proxy/codec"
)

// StreamDirector returns a gRPC ClientConn to be used to forward the call to.
//
// The presence of the `Context` allows for rich filtering, e.g. based on Metadata (headers).
// If no handling is meant to be done, a `codes.NotImplemented` gRPC error should be returned.
//
// The context returned from this function should be the context for the *outgoing* (to backend) call. In case you want
// to forward any Metadata between the inbound request and outbound requests, you should do it manually. However, you
// *must* propagate the cancel function (`context.WithCancel`) of the inbound context to the one returned.
//
// It is worth noting that the StreamDirector will be fired *after* all server-side stream interceptors
// are invoked. So decisions around authorization, monitoring etc. are better to be handled there.
//
// See the rather rich example.
type StreamDirector func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, func(), error)

type grpcManager struct {
	lock  sync.Mutex
	addr  string
	conns map[string]*grpc.ClientConn
}

func NewManager(addr string) *grpcManager {
	return &grpcManager{conns: make(map[string]*grpc.ClientConn), addr: addr}
}

func (m *grpcManager) StreamDirector(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, func(), error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// todo: support addr metadata mapping
	conn, ok := m.conns[m.addr]
	if ok {
		return ctx, conn, func() {}, nil
	}

	fmt.Println("create conn")

	conn, err := grpc.DialContext(ctx, m.addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.CallContentSubtype((&codec.Proxy{}).Name())))
	if err != nil {
		return ctx, nil, func() {}, err
	}

	m.conns[m.addr] = conn
	return ctx, conn, func() {}, nil
}
