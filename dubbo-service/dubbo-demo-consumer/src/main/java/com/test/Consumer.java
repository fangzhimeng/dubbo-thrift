package com.test;

import com.test.thrift.Demo;
import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.text.SimpleDateFormat;
import java.util.Date;

public class Consumer {
    public static void main(String[] args) throws Exception {
        ClassPathXmlApplicationContext context =
                new ClassPathXmlApplicationContext("dubbo-demo-consumer.xml");
        context.start();
        Demo.Iface demo = (Demo.Iface) context.getBean("demoService");
        for (int i = 0; i < Integer.MAX_VALUE; i ++) {
            try {
                String hello = demo.sayHello("world" + i);
                System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " + hello);
            } catch (Exception e) {
                e.printStackTrace();
            }
            Thread.sleep(2000);
        }
        context.close();
    }

}
