package com.grab.tracking.service;

import lombok.extern.slf4j.Slf4j;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

import java.time.Duration;
import java.util.UUID;

@Slf4j
@Service
public class ActiveOrderTracker {

    private final ReactiveRedisTemplate<String, String> stringRedisTemplate;
    private static final String ORDER_DRIVER_KEY = "order:driver:";
    private static final Duration ORDER_TTL = Duration.ofHours(24);

    public ActiveOrderTracker(ReactiveRedisTemplate<String, String> stringRedisTemplate) {
        this.stringRedisTemplate = stringRedisTemplate;
    }

    public Mono<Boolean> trackOrderDriver(UUID orderId, UUID driverId) {
        return stringRedisTemplate.opsForValue()
                .set(ORDER_DRIVER_KEY + orderId, driverId.toString(), ORDER_TTL)
                .doOnSuccess(ok -> log.info("Tracking order={} driver={}", orderId, driverId));
    }

    public Mono<String> getDriverForOrder(UUID orderId) {
        return stringRedisTemplate.opsForValue().get(ORDER_DRIVER_KEY + orderId);
    }

    public Mono<Boolean> stopTracking(UUID orderId) {
        return stringRedisTemplate.delete(ORDER_DRIVER_KEY + orderId)
                .map(count -> count > 0);
    }
}
