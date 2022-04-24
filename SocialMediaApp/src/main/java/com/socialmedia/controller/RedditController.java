package com.socialmedia.controller;

import com.socialmedia.service.SocialMediaService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class RedditController {
    @Autowired
    private SocialMediaService socialMediaService;

    // Save operation
    @PostMapping("/reddit")
    public HttpStatus saveDepartment(@RequestParam(value="folder") String folder)
    {
        return socialMediaService.saveSocialMedia(folder);
    }

}
