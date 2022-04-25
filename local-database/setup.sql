CREATE DATABASE IF NOT EXISTS socialmedia CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS reddit (
    date_posted TIMESTAMP,
    region VARCHAR(50),
    subreddit VARCHAR(20),
    post_title VARCHAR(255),
    upvotes INTEGER,
    date_stored TIMESTAMP,
    comments INTEGER,
    PRIMARY KEY (date_posted, post_title, upvotes)
);

CREATE TABLE IF NOT EXISTS twitter (
    date_posted TIMESTAMP,
    region VARCHAR(50),
    tweet VARCHAR(280),
    likes INTEGER,
    retweets INTEGER,
    date_stored TIMESTAMP,
    comments INTEGER,
    PRIMARY KEY (date_posted, tweet)
);

-- INSERT INTO reddit ( date_posted, region, subreddit, post_title, upvotes, date_stored, comments) VALUES ( '2011-08-21 14:11:09', 'US', 'worldnews', 'test post title 1', '1', '2011-08-21 14:11:09', '10' );
-- INSERT INTO reddit ( date_posted, region, subreddit, post_title, upvotes, date_stored, comments) VALUES ( '2011-08-22 14:11:09', 'US', 'news', 'test post title 2', '2', '2011-08-21 14:11:09', '20' );
-- INSERT INTO reddit ( date_posted, region, subreddit, post_title, upvotes, date_stored, comments) VALUES ( '2011-08-23 14:11:09', 'US', 'world', 'test post title 3', '3', '2011-08-21 14:11:09', '30' );

INSERT INTO twitter ( date_posted, region, tweet, likes, retweets, date_stored, comments) VALUES ( '2011-08-21 14:11:09', 'US', 'test tweet 1', '10', '100', '2011-08-21 14:11:09', '10' );
INSERT INTO twitter ( date_posted, region, tweet, likes, retweets, date_stored, comments) VALUES ( '2011-08-22 14:11:09', 'US', 'test tweet 2', '20', '200', '2011-08-21 14:11:09', '11' );
INSERT INTO twitter ( date_posted, region, tweet, likes, retweets, date_stored, comments) VALUES ( '2011-08-23 14:11:09', 'US', 'test tweet 3', '30', '300', '2011-08-21 14:11:09', '12' );