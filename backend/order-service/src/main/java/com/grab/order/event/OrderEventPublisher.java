package com.grab.order.event;

import com.grab.order.model.Order;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

import java.time.LocalDateTime;

@Slf4j
@Component
@RequiredArgsConstructor
public class OrderEventPublisher {

    private final KafkaTemplate<String, OrderEvent> kafkaTemplate;

    @Value("${kafka.topics.order-created}")
    private String orderCreatedTopic;

    @Value("${kafka.topics.order-accepted}")
    private String orderAcceptedTopic;

    @Value("${kafka.topics.order-completed}")
    private String orderCompletedTopic;

    @Value("${kafka.topics.order-cancelled}")
    private String orderCancelledTopic;

    public void publishOrderCreated(Order order) {
        publish(orderCreatedTopic, buildEvent(order, OrderEvent.EventType.ORDER_CREATED));
    }

    public void publishOrderAccepted(Order order) {
        publish(orderAcceptedTopic, buildEvent(order, OrderEvent.EventType.ORDER_ACCEPTED));
    }

    public void publishOrderDelivered(Order order) {
        publish(orderCompletedTopic, buildEvent(order, OrderEvent.EventType.ORDER_DELIVERED));
    }

    public void publishOrderCancelled(Order order) {
        publish(orderCancelledTopic, buildEvent(order, OrderEvent.EventType.ORDER_CANCELLED));
    }

    private void publish(String topic, OrderEvent event) {
        kafkaTemplate.send(topic, event.getOrderId().toString(), event)
                .whenComplete((result, ex) -> {
                    if (ex != null) {
                        log.error("Failed to publish event {} to topic {}: {}", event.getEventType(), topic, ex.getMessage());
                    } else {
                        log.info("Published event {} for orderId={}", event.getEventType(), event.getOrderId());
                    }
                });
    }

    private OrderEvent buildEvent(Order order, OrderEvent.EventType type) {
        return OrderEvent.builder()
                .eventType(type.name())
                .orderId(order.getId())
                .customerId(order.getCustomerId())
                .driverId(order.getDriverId())
                .status(order.getStatus().name())
                .pickupAddress(order.getPickupAddress())
                .dropoffAddress(order.getDropoffAddress())
                .totalPrice(order.getTotalPrice())
                .occurredAt(LocalDateTime.now())
                .build();
    }
}
