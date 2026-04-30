package com.grab.tracking.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class DriverLocation {
    private UUID driverId;
    private double latitude;
    private double longitude;
    private double heading;
    private double speed;
    private Instant updatedAt;
}
