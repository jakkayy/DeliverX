package com.grab.tracking.kafka;

import com.grab.tracking.service.ActiveOrderTracker;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import java.util.Map;
import java.util.UUID;

@Slf4j
@Component
@RequiredArgsConstructor
public class OrderEventConsumer {

    private final ActiveOrderTracker activeOrderTracker;

    @KafkaListener(topics = "${kafka.topics.order-accepted}", groupId = "tracking-service")
    public void onOrderAccepted(Map<String, Object> event) {
        try {
            String orderId = event.get("orderId").toString();
            String driverId = event.get("driverId").toString();
            activeOrderTracker.trackOrderDriver(
                    UUID.fromString(orderId),
                    UUID.fromString(driverId)
            ).subscribe();
            log.info("Tracking started: orderId={} driverId={}", orderId, driverId);
        } catch (Exception e) {
            log.error("Failed to process order-accepted event: {}", e.getMessage());
        }
    }

    @KafkaListener(topics = "${kafka.topics.order-completed}", groupId = "tracking-service")
    public void onOrderCompleted(Map<String, Object> event) {
        try {
            String orderId = event.get("orderId").toString();
            activeOrderTracker.stopTracking(UUID.fromString(orderId)).subscribe();
            log.info("Tracking stopped: orderId={}", orderId);
        } catch (Exception e) {
            log.error("Failed to process order-completed event: {}", e.getMessage());
        }
    }

    @KafkaListener(topics = "${kafka.topics.order-cancelled}", groupId = "tracking-service")
    public void onOrderCancelled(Map<String, Object> event) {
        try {
            String orderId = event.get("orderId").toString();
            activeOrderTracker.stopTracking(UUID.fromString(orderId)).subscribe();
            log.info("Tracking stopped (cancelled): orderId={}", orderId);
        } catch (Exception e) {
            log.error("Failed to process order-cancelled event: {}", e.getMessage());
        }
    }
}
