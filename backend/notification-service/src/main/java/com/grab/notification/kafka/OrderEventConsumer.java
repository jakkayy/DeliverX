package com.grab.notification.kafka;

import com.grab.notification.model.Notification.NotificationType;
import com.grab.notification.service.NotificationService;
import com.grab.notification.service.UserFcmTokenResolver;
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

    private final NotificationService notificationService;
    private final UserFcmTokenResolver tokenResolver;

    @KafkaListener(topics = "${kafka.topics.order-created}", groupId = "notification-service")
    public void onOrderCreated(Map<String, Object> event) {
        try {
            UUID customerId = UUID.fromString(event.get("customerId").toString());
            String orderId = event.get("orderId").toString();

            String fcmToken = tokenResolver.getToken(customerId);
            notificationService.sendAndSave(
                    customerId, fcmToken,
                    "Order Placed", "Your order has been placed. Looking for a driver...",
                    NotificationType.ORDER_CREATED, orderId,
                    Map.of("orderId", orderId, "screen", "order_detail")
            );
        } catch (Exception e) {
            log.error("Failed to process order-created event: {}", e.getMessage());
        }
    }

    @KafkaListener(topics = "${kafka.topics.order-accepted}", groupId = "notification-service")
    public void onOrderAccepted(Map<String, Object> event) {
        try {
            UUID customerId = UUID.fromString(event.get("customerId").toString());
            String orderId = event.get("orderId").toString();

            String fcmToken = tokenResolver.getToken(customerId);
            notificationService.sendAndSave(
                    customerId, fcmToken,
                    "Driver Found!", "A driver has accepted your order and is on the way.",
                    NotificationType.ORDER_ACCEPTED, orderId,
                    Map.of("orderId", orderId, "screen", "tracking")
            );
        } catch (Exception e) {
            log.error("Failed to process order-accepted event: {}", e.getMessage());
        }
    }

    @KafkaListener(topics = "${kafka.topics.order-completed}", groupId = "notification-service")
    public void onOrderDelivered(Map<String, Object> event) {
        try {
            UUID customerId = UUID.fromString(event.get("customerId").toString());
            String orderId = event.get("orderId").toString();

            String fcmToken = tokenResolver.getToken(customerId);
            notificationService.sendAndSave(
                    customerId, fcmToken,
                    "Order Delivered!", "Your order has been delivered successfully.",
                    NotificationType.ORDER_DELIVERED, orderId,
                    Map.of("orderId", orderId, "screen", "order_detail")
            );
        } catch (Exception e) {
            log.error("Failed to process order-completed event: {}", e.getMessage());
        }
    }

    @KafkaListener(topics = "${kafka.topics.order-cancelled}", groupId = "notification-service")
    public void onOrderCancelled(Map<String, Object> event) {
        try {
            UUID customerId = UUID.fromString(event.get("customerId").toString());
            String orderId = event.get("orderId").toString();
            UUID driverId = event.get("driverId") != null
                    ? UUID.fromString(event.get("driverId").toString()) : null;

            String customerToken = tokenResolver.getToken(customerId);
            notificationService.sendAndSave(
                    customerId, customerToken,
                    "Order Cancelled", "Your order has been cancelled.",
                    NotificationType.ORDER_CANCELLED, orderId,
                    Map.of("orderId", orderId, "screen", "home")
            );

            if (driverId != null) {
                String driverToken = tokenResolver.getToken(driverId);
                notificationService.sendAndSave(
                        driverId, driverToken,
                        "Order Cancelled", "The order has been cancelled by the customer.",
                        NotificationType.ORDER_CANCELLED, orderId,
                        Map.of("orderId", orderId, "screen", "driver_home")
                );
            }
        } catch (Exception e) {
            log.error("Failed to process order-cancelled event: {}", e.getMessage());
        }
    }
}
