package com.grab.notification.service;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;

import java.util.Map;
import java.util.UUID;

@Slf4j
@Service
public class UserFcmTokenResolver {

    private final StringRedisTemplate redisTemplate;
    private final RestClient restClient;

    private static final String FCM_TOKEN_PREFIX = "user:fcm:";

    public UserFcmTokenResolver(StringRedisTemplate redisTemplate,
                                 @Value("${services.user-service.url:http://user-service:8082}") String userServiceUrl) {
        this.redisTemplate = redisTemplate;
        this.restClient = RestClient.builder().baseUrl(userServiceUrl).build();
    }

    public String getToken(UUID userId) {
        String cached = redisTemplate.opsForValue().get(FCM_TOKEN_PREFIX + userId);
        if (cached != null) {
            return cached;
        }

        try {
            Map<?, ?> response = restClient.get()
                    .uri("/api/v1/users/{id}", userId)
                    .retrieve()
                    .body(Map.class);

            if (response != null && response.get("fcmToken") != null) {
                String token = response.get("fcmToken").toString();
                redisTemplate.opsForValue().set(FCM_TOKEN_PREFIX + userId, token);
                return token;
            }
        } catch (Exception e) {
            log.warn("Could not fetch FCM token for userId={}: {}", userId, e.getMessage());
        }

        return null;
    }

    public void cacheToken(UUID userId, String fcmToken) {
        redisTemplate.opsForValue().set(FCM_TOKEN_PREFIX + userId, fcmToken);
    }
}
