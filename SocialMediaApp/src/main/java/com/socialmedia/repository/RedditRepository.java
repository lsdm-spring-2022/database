package com.socialmedia.repository;

import com.socialmedia.entity.Reddit;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface RedditRepository extends CrudRepository<Reddit, Long> {
}
