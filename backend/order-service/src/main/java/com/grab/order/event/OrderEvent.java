package com.grab.order.event;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OrderEvent {

    private String eventType;
    private UUID orderId;
    private UUID customerId;
    private UUID driverId;
    private String status;
    private String pickupAddress;
    private String dropoffAddress;
    private BigDecimal totalPrice;
    private LocalDateTime occurredAt;

    public enum EventType {
        ORDER_CREATED,
        ORDER_ACCEPTED,
        ORDER_PICKUP,
        ORDER_IN_TRANSIT,
        ORDER_DELIVERED,
        ORDER_CANCELLED
    }
}
