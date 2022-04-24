package com.socialmedia.entity;

import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;

import javax.persistence.Entity;
import javax.persistence.Id;
import java.sql.Timestamp;

@Entity
@AllArgsConstructor
@NoArgsConstructor
public class Reddit {
    @Id
    private RedditId redditId;

    private String region;
    private int upvotes;
    private String post_title;
    private Timestamp date_stored;
    private int comments;
}
