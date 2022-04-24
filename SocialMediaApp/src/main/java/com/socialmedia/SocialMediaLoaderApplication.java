package com.socialmedia;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class SocialMediaLoaderApplication{

    public static void main(String[] args) {
        System.out.println("Starting The Application");
        SpringApplication.run(SocialMediaLoaderApplication.class, args);
    }
}
