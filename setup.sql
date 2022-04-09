CREATE DATABASE IF NOT EXISTS socialmedia;

CREATE TABLE IF NOT EXISTS reddit (
    date_posted TIMESTAMP,
    region VARCHAR(50),
    subreddit VARCHAR(20),
    post_title VARCHAR(255),
    upvotes INTEGER(11),
    date_stored TIMESTAMP,
    comments INTEGER(11),
    PRIMARY KEY (date_posted, subreddit)
);

INSERT INTO reddit ( date_posted, region, subreddit, post_title, upvotes, date_stored, comments) VALUES ( '2011-08-21 14:11:09', 'US', 'worldnews', 'test post title 1', '1', '2011-08-21 14:11:09', '10' );
INSERT INTO reddit ( date_posted, region, subreddit, post_title, upvotes, date_stored, comments) VALUES ( '2011-08-22 14:11:09', 'US', 'news', 'test post title 2', '2', '2011-08-21 14:11:09', '20' );
INSERT INTO reddit ( date_posted, region, subreddit, post_title, upvotes, date_stored, comments) VALUES ( '2011-08-23 14:11:09', 'US', 'world', 'test post title 3', '3', '2011-08-21 14:11:09', '30' );