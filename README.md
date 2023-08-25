# grpc-go-beyond
<!--
[![Go Report Card](https://goreportcard.com/badge/github.com/zcong1993/template-go-cli)](https://goreportcard.com/report/github.com/zcong1993/template-go-cli)
-->

> grpc schemaless proxy

```bash
# start bizhello
go run ./cmd/bizhello
# start biztest
go run ./cmd/biztest
# start proxy
go run ./cmd/proxy
# run test
# 1. biz1 call Echo
go run ./cmd/client -method Echo -sign-key biz1 -app biz1
# 2. biz1 call Test
go run ./cmd/client -method Test -sign-key biz1 -app biz1
# 3. biz2 call Echo2
go run ./cmd/client -method Echo2 -sign-key biz2 -app biz2
# 4. no access
go run ./cmd/client -method Echo2 -sign-key biz1 -app biz1
# 5. sign error
go run ./cmd/client -method Echo -sign-key biz2 -app biz1
```

## License

MIT &copy; zcong1993
