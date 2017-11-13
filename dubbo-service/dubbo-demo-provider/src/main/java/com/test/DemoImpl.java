package com.test;
import com.alibaba.dubbo.rpc.RpcContext;
import com.test.thrift.Demo;
import org.apache.thrift.TException;

import java.text.SimpleDateFormat;
import java.util.Date;

public class DemoImpl implements Demo.Iface {
    public String sayHello(String hello) throws TException {
        System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] Hello " + hello + ", request from consumer: " + RpcContext.getContext().getRemoteAddress());
        return "Hello " + hello + ", response form provider: " + RpcContext.getContext().getLocalAddress();
    }
}
