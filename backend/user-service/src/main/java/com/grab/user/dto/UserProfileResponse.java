package com.grab.user.dto;

import lombok.Builder;
import lombok.Data;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Builder
public class UserProfileResponse {
    private UUID id;
    private String name;
    private String phone;
    private String email;
    private String role;
    private String profileImage;
    private Boolean isActive;
    private LocalDateTime createdAt;
}
