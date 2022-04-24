package com.socialmedia.service;

import org.springframework.http.HttpStatus;

public interface SocialMediaService {
    HttpStatus saveSocialMedia(String location);
}
