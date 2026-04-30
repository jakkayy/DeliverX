package com.grab.user.controller;

import com.grab.user.dto.UpdateProfileRequest;
import com.grab.user.dto.UserProfileResponse;
import com.grab.user.service.UserService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping("/me")
    public ResponseEntity<UserProfileResponse> getMyProfile(
            @RequestHeader("X-User-Id") String userId) {
        return ResponseEntity.ok(userService.getProfile(UUID.fromString(userId)));
    }

    @GetMapping("/{id}")
    public ResponseEntity<UserProfileResponse> getProfile(@PathVariable UUID id) {
        return ResponseEntity.ok(userService.getProfile(id));
    }

    @PutMapping("/me")
    public ResponseEntity<UserProfileResponse> updateProfile(
            @RequestHeader("X-User-Id") String userId,
            @Valid @RequestBody UpdateProfileRequest request) {
        return ResponseEntity.ok(userService.updateProfile(UUID.fromString(userId), request));
    }

    @DeleteMapping("/me")
    public ResponseEntity<Void> deactivate(@RequestHeader("X-User-Id") String userId) {
        userService.deactivateAccount(UUID.fromString(userId));
        return ResponseEntity.noContent().build();
    }
}
