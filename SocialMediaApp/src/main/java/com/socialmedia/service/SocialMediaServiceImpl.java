package com.socialmedia.service;

import com.socialmedia.helper.CSVHelper;
import com.socialmedia.entity.Reddit;
import com.socialmedia.repository.RedditRepository;
import org.apache.commons.io.IOUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.mock.web.MockMultipartFile;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
public class SocialMediaServiceImpl implements SocialMediaService{
    @Autowired
    RedditRepository repository;

    @Override
    public HttpStatus saveSocialMedia(String folder) {
        Path path = Paths.get("/Users/dezerelgraham/Downloads/" + folder);
        System.out.println(path);
        List<Path> paths = null;
        try {
            paths = listFiles(path);
        } catch (IOException e) {
            e.printStackTrace();
        }
        assert paths != null;
        paths.forEach(x -> {
            try {
                System.out.println(x);
                storeData(x);
            } catch (IOException e) {
                e.printStackTrace();
            }
        });
        return HttpStatus.CREATED;
    }

    private List<Path> listFiles(Path path) throws IOException {
        List<Path> result;
        try (Stream<Path> walk = Files.walk(path)) {
            result = walk.filter(Files::isRegularFile)
                    .collect(Collectors.toList());
        }
        return result;
    }

    private void storeData(Path x) throws IOException {
        File f = new File(x.toString());
        FileInputStream input = new FileInputStream(f);
        MultipartFile file = new MockMultipartFile("file",
                f.getName(), "text/plain", IOUtils.toByteArray(input));
        CSVHelper csvHelper = new CSVHelper();
        List<Reddit> subreddits = csvHelper.csvToReddit(file.getInputStream());
        AtomicInteger errorcount = new AtomicInteger();
        subreddits.forEach(reddit -> {
            try{
                repository.save(reddit);
            }catch(Exception e){
                errorcount.getAndIncrement();
            }
        });
        System.out.println(errorcount.get());
    }
}
