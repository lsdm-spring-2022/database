package com.socialmedia.helper;

import com.socialmedia.entity.Reddit;
import com.socialmedia.entity.RedditId;
import org.apache.commons.csv.CSVFormat;
import org.apache.commons.csv.CSVParser;
import org.apache.commons.csv.CSVRecord;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.sql.Timestamp;
import java.time.Instant;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

public class CSVHelper {
    public List<Reddit> csvToReddit(InputStream is) {
        try (BufferedReader fileReader = new BufferedReader(new InputStreamReader(is, "UTF-8"));
             CSVParser csvParser = new CSVParser(fileReader,
                     CSVFormat.DEFAULT.withFirstRecordAsHeader().withIgnoreHeaderCase().withAllowDuplicateHeaderNames()
                             .withAllowMissingColumnNames().withQuote(null).withTrim());) {
            List<Reddit> subreddits = new ArrayList<>();
            Iterable<CSVRecord> csvRecords = csvParser.getRecords();

            for (CSVRecord csvRecord : csvRecords) {
                RedditId redditId = new RedditId(covertToTimestamp(csvRecord.get("created_utc")),
                        csvRecord.get("subreddit"));
                Reddit reddit = new Reddit(
                        redditId,
                        csvRecord.get("subreddit"),
                        Integer.parseInt(csvRecord.get("score")),
                        csvRecord.get("title"),
                        Timestamp.valueOf(LocalDateTime.now()),
                        Integer.parseInt(csvRecord.get("num_comments"))
                );
                subreddits.add(reddit);
            }
            return subreddits;
        } catch (IOException e) {
            e.printStackTrace();
            throw new RuntimeException("fail to parse CSV file: " + e.getMessage());
        }
    }

    private Timestamp covertToTimestamp(String created_utc) {
        Timestamp created = Timestamp.from(Instant.now());
        try{
            created = Timestamp.from(Instant.ofEpochSecond(Long.parseLong(created_utc)));
        }catch(Exception e){
            e.printStackTrace();
        }
        return created;
    }
}
