CREATE DATABASE IF NOT EXISTS socialmedia CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS reddit (
    date_posted TIMESTAMP,
    region VARCHAR(50),
    subreddit VARCHAR(100),
    post_title VARCHAR(255),
    upvotes INTEGER,
    date_stored TIMESTAMP,
    comments INTEGER,
    PRIMARY KEY (date_posted, post_title, upvotes)
);

CREATE TABLE IF NOT EXISTS twitter (
    country VARCHAR(50),
    created_at TIMESTAMP,
    trend VARCHAR(100),
    tweet_volume INTEGER,
    as_of TIMESTAMP,
    date_retrieved TIMESTAMP,
    PRIMARY KEY (trend, country, tweet_volume)
);
