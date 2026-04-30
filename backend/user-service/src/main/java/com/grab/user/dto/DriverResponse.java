package com.grab.user.dto;

import lombok.Builder;
import lombok.Data;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Builder
public class DriverResponse {
    private UUID id;
    private UUID userId;
    private String name;
    private String phone;
    private String profileImage;
    private String vehicleType;
    private String licensePlate;
    private String vehicleModel;
    private BigDecimal rating;
    private Integer totalTrips;
    private Boolean isAvailable;
    private LocalDateTime createdAt;
}
