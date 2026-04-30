package com.grab.user.controller;

import com.grab.user.dto.DriverResponse;
import com.grab.user.dto.RegisterDriverRequest;
import com.grab.user.service.DriverService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/drivers")
@RequiredArgsConstructor
public class DriverController {

    private final DriverService driverService;

    @PostMapping("/register")
    public ResponseEntity<DriverResponse> register(
            @RequestHeader("X-User-Id") String userId,
            @Valid @RequestBody RegisterDriverRequest request) {
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(driverService.registerDriver(UUID.fromString(userId), request));
    }

    @GetMapping("/me")
    public ResponseEntity<DriverResponse> getMyProfile(
            @RequestHeader("X-User-Id") String userId) {
        return ResponseEntity.ok(driverService.getDriverByUserId(UUID.fromString(userId)));
    }

    @GetMapping("/{id}")
    public ResponseEntity<DriverResponse> getDriver(@PathVariable UUID id) {
        return ResponseEntity.ok(driverService.getDriverById(id));
    }

    @PatchMapping("/me/availability")
    public ResponseEntity<DriverResponse> updateAvailability(
            @RequestHeader("X-User-Id") String userId,
            @RequestBody Map<String, Boolean> body) {
        boolean isAvailable = body.getOrDefault("isAvailable", false);
        return ResponseEntity.ok(driverService.updateAvailability(UUID.fromString(userId), isAvailable));
    }

    @GetMapping("/available")
    public ResponseEntity<List<DriverResponse>> getAvailableDrivers() {
        return ResponseEntity.ok(driverService.getAvailableDrivers());
    }
}
