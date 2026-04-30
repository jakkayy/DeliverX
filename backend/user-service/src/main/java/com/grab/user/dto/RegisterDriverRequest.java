package com.grab.user.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

@Data
public class RegisterDriverRequest {

    @NotNull(message = "Vehicle type is required")
    private String vehicleType;

    @NotBlank(message = "License plate is required")
    private String licensePlate;

    private String vehicleModel;
}
