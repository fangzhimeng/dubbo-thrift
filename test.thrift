namespace java com.test.thrift
namespace go test

service Demo {
    string sayHello( 1:required string hello );
}