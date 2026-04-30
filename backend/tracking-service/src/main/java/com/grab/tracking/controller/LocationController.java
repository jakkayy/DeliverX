package com.grab.tracking.controller;

import com.grab.tracking.model.DriverLocation;
import com.grab.tracking.model.LocationUpdateRequest;
import com.grab.tracking.service.ActiveOrderTracker;
import com.grab.tracking.service.LocationService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.time.Duration;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/tracking")
@RequiredArgsConstructor
public class LocationController {

    private final LocationService locationService;
    private final ActiveOrderTracker activeOrderTracker;

    @PostMapping("/location")
    public Mono<ResponseEntity<DriverLocation>> updateLocation(
            @RequestHeader("X-User-Id") String driverId,
            @Valid @RequestBody LocationUpdateRequest request) {
        return locationService.updateLocation(UUID.fromString(driverId), request)
                .map(ResponseEntity::ok);
    }

    @GetMapping("/driver/{driverId}")
    public Mono<ResponseEntity<DriverLocation>> getDriverLocation(@PathVariable UUID driverId) {
        return locationService.getDriverLocation(driverId)
                .map(ResponseEntity::ok)
                .defaultIfEmpty(ResponseEntity.notFound().build());
    }

    @GetMapping(value = "/order/{orderId}/stream", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<String> streamOrderTracking(@PathVariable String orderId) {
        return locationService.subscribeToOrderTracking(orderId)
                .timeout(Duration.ofMinutes(30));
    }

    @PostMapping("/order/{orderId}/publish")
    public Mono<ResponseEntity<Void>> publishLocationToOrder(
            @RequestHeader("X-User-Id") String driverId,
            @PathVariable String orderId,
            @Valid @RequestBody LocationUpdateRequest request) {
        return locationService.updateLocation(UUID.fromString(driverId), request)
                .flatMap(location -> locationService.publishLocationToOrder(orderId, location))
                .map(count -> ResponseEntity.<Void>ok().build());
    }

    @GetMapping("/order/{orderId}/driver")
    public Mono<ResponseEntity<String>> getOrderDriver(@PathVariable String orderId) {
        return activeOrderTracker.getDriverForOrder(UUID.fromString(orderId))
                .map(ResponseEntity::ok)
                .defaultIfEmpty(ResponseEntity.notFound().build());
    }
}
