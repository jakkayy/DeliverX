package com.grab.tracking.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.grab.tracking.model.DriverLocation;
import com.grab.tracking.model.LocationUpdateRequest;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.time.Duration;
import java.time.Instant;
import java.util.UUID;

@Slf4j
@Service
@RequiredArgsConstructor
public class LocationService {

    private final ReactiveRedisTemplate<String, DriverLocation> locationRedisTemplate;
    private final ReactiveRedisTemplate<String, String> stringRedisTemplate;
    private final ObjectMapper objectMapper;

    @Value("${redis.key.driver-location-prefix:driver:location:}")
    private String locationKeyPrefix;

    @Value("${redis.key.order-channel-prefix:order:channel:}")
    private String channelKeyPrefix;

    private static final Duration LOCATION_TTL = Duration.ofMinutes(10);

    public Mono<DriverLocation> updateLocation(UUID driverId, LocationUpdateRequest request) {
        DriverLocation location = DriverLocation.builder()
                .driverId(driverId)
                .latitude(request.getLatitude())
                .longitude(request.getLongitude())
                .heading(request.getHeading())
                .speed(request.getSpeed())
                .updatedAt(Instant.now())
                .build();

        String key = locationKeyPrefix + driverId;

        return locationRedisTemplate.opsForValue()
                .set(key, location, LOCATION_TTL)
                .then(publishLocationUpdate(driverId, location))
                .thenReturn(location)
                .doOnSuccess(l -> log.debug("Location updated: driver={} lat={} lng={}", driverId, l.getLatitude(), l.getLongitude()));
    }

    public Mono<DriverLocation> getDriverLocation(UUID driverId) {
        return locationRedisTemplate.opsForValue()
                .get(locationKeyPrefix + driverId);
    }

    public Flux<String> subscribeToOrderTracking(String orderId) {
        String channel = channelKeyPrefix + orderId;
        return stringRedisTemplate.listenToChannel(channel)
                .map(message -> message.getMessage())
                .doOnSubscribe(s -> log.info("Client subscribed to order channel: {}", orderId))
                .doOnCancel(() -> log.info("Client unsubscribed from order channel: {}", orderId));
    }

    public Mono<Long> publishLocationToOrder(String orderId, DriverLocation location) {
        String channel = channelKeyPrefix + orderId;
        try {
            String message = objectMapper.writeValueAsString(location);
            return stringRedisTemplate.convertAndSend(channel, message);
        } catch (JsonProcessingException e) {
            return Mono.error(e);
        }
    }

    private Mono<Long> publishLocationUpdate(UUID driverId, DriverLocation location) {
        String channel = "driver:updates:" + driverId;
        try {
            String message = objectMapper.writeValueAsString(location);
            return stringRedisTemplate.convertAndSend(channel, message);
        } catch (JsonProcessingException e) {
            return Mono.error(e);
        }
    }
}
