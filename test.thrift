namespace java com.alibaba.dubbo.rpc.thrift
namespace go test

service Demo {
    string sayHello( 1:required string hello );
}