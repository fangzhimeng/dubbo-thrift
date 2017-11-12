package dubbo_thrift

import (
	"dubbo-thrift/test"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"testing"
	"time"
)

const (
	serviceAddr = "192.168.1.108:9000"
)

func TestClient(t *testing.T) {
	client, err := getThriftClient(serviceAddr, false)
	if err != nil {
		t.Fatalf("set up client error %s", err)
	}
	for i := 0; i < 100; i++ {
		req := "hello"
		fmt.Println("send " + req)
		res, err := client.SayHello(req)
		if err != nil {
			t.Errorf("read client error %s", err)
		}
		fmt.Println("received " + res)
	}
}

func getThriftClient(addr string, dubboTransport bool) (test.Demo, error) {
	socket, err := thrift.NewTSocketTimeout(addr, 1000*time.Millisecond)
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return nil, err
	}
	protocolFactory := NewTDubboProtocolFactory("com.test.thrift.Demo$Iface")
	var transportFactory thrift.TTransportFactory
	if dubboTransport {
		transportFactory = NewTDubboTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	}
	transport := transportFactory.GetTransport(socket)
	if err := transport.Open(); err != nil {
		return nil, err
	}
	return test.NewDemoClientFactory(transport, protocolFactory), nil
}

func getThriftServer(handler test.Demo, addr string, dubboTransport bool) (thrift.TServer, error) {
	protocolFactory := NewTDubboProtocolFactory("com.test.thrift.Demo$Iface")
	var transportFactory thrift.TTransportFactory
	if dubboTransport {
		transportFactory = NewTDubboTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	}
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return nil, err
	}
	processor := test.NewDemoProcessor(handler)
	return thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory), nil
}
