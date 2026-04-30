package com.grab.order.dto;

import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

import java.math.BigDecimal;
import java.util.List;

@Data
public class CreateOrderRequest {

    @NotBlank(message = "Pickup address is required")
    private String pickupAddress;

    @NotNull(message = "Pickup latitude is required")
    private BigDecimal pickupLat;

    @NotNull(message = "Pickup longitude is required")
    private BigDecimal pickupLng;

    @NotBlank(message = "Dropoff address is required")
    private String dropoffAddress;

    @NotNull(message = "Dropoff latitude is required")
    private BigDecimal dropoffLat;

    @NotNull(message = "Dropoff longitude is required")
    private BigDecimal dropoffLng;

    private String note;

    @Valid
    private List<OrderItemRequest> items;

    @Data
    public static class OrderItemRequest {
        @NotBlank
        private String name;
        private String description;
        private Integer quantity = 1;
        private BigDecimal weightKg;
        private String imageUrl;
    }
}
