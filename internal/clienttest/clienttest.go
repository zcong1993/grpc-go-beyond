package clienttest

type Tester interface {
	TestEcho()
	TestServerStream()
	TestClientStream()
	TestDuplexStream()
}
