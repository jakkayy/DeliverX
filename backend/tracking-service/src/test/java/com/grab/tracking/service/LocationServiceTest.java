package com.grab.tracking.service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.grab.tracking.model.DriverLocation;
import com.grab.tracking.model.LocationUpdateRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.data.redis.core.ReactiveValueOperations;
import reactor.core.publisher.Mono;
import reactor.test.StepVerifier;

import java.util.UUID;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class LocationServiceTest {

    @Mock
    private ReactiveRedisTemplate<String, DriverLocation> locationRedisTemplate;

    @Mock
    private ReactiveRedisTemplate<String, String> stringRedisTemplate;

    @Mock
    private ReactiveValueOperations<String, DriverLocation> locationOps;

    @Mock
    private ReactiveValueOperations<String, String> stringOps;

    private LocationService locationService;
    private UUID driverId;

    @BeforeEach
    void setUp() {
        ObjectMapper objectMapper = new ObjectMapper()
                .registerModule(new JavaTimeModule())
                .disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);

        locationService = new LocationService(locationRedisTemplate, stringRedisTemplate, objectMapper);
        driverId = UUID.randomUUID();
    }

    @Test
    @DisplayName("Update location stores in Redis and publishes event")
    void updateLocation_success() {
        LocationUpdateRequest request = new LocationUpdateRequest();
        request.setLatitude(13.7563);
        request.setLongitude(100.5018);
        request.setHeading(90.0);
        request.setSpeed(30.0);

        when(locationRedisTemplate.opsForValue()).thenReturn(locationOps);
        when(locationOps.set(anyString(), any(DriverLocation.class), any())).thenReturn(Mono.just(true));
        when(stringRedisTemplate.convertAndSend(anyString(), anyString())).thenReturn(Mono.just(1L));

        StepVerifier.create(locationService.updateLocation(driverId, request))
                .expectNextMatches(loc ->
                        loc.getDriverId().equals(driverId) &&
                        loc.getLatitude() == 13.7563 &&
                        loc.getLongitude() == 100.5018)
                .verifyComplete();
    }

    @Test
    @DisplayName("Get driver location returns empty when not found")
    void getDriverLocation_notFound() {
        when(locationRedisTemplate.opsForValue()).thenReturn(locationOps);
        when(locationOps.get(anyString())).thenReturn(Mono.empty());

        StepVerifier.create(locationService.getDriverLocation(driverId))
                .verifyComplete();
    }
}
