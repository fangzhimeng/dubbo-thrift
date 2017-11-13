package com.test;

import org.springframework.context.support.ClassPathXmlApplicationContext;

public class Provider {
    public static void main(String[] args) throws Exception {
        ClassPathXmlApplicationContext context =
                new ClassPathXmlApplicationContext("dubbo-demo-provider.xml");
        context.start();
        System.out.println("context started");
        System.in.read();
    }
}
