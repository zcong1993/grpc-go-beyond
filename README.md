# grpc-go-beyond
<!--
[![Go Report Card](https://goreportcard.com/badge/github.com/zcong1993/template-go-cli)](https://goreportcard.com/report/github.com/zcong1993/template-go-cli)
-->

> low level grpc go server/client demo.

- [internal/client/stream.go](internal/client/stream.go) 使用 `grpc ClientStream` 实现带有类型的 grpc client
- [internal/server/stream.go](internal/server/stream.go) 使用 `grpc StreamHandler` 实现 grpc server
- [internal/clienttest/raw.go](internal/clienttest/raw.go) 使用 `grpc ClientStream` 直接和服务端交互

## License

MIT &copy; zcong1993
