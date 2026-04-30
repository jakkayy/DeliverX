package com.grab.order.service;

import com.grab.order.dto.CreateOrderRequest;
import com.grab.order.dto.OrderResponse;
import com.grab.order.event.OrderEventPublisher;
import com.grab.order.exception.OrderNotFoundException;
import com.grab.order.exception.OrderStatusException;
import com.grab.order.model.Order;
import com.grab.order.model.Order.OrderStatus;
import com.grab.order.model.OrderItem;
import com.grab.order.repository.OrderRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class OrderService {

    private final OrderRepository orderRepository;
    private final OrderEventPublisher eventPublisher;
    private final PricingService pricingService;

    @Transactional
    public OrderResponse createOrder(UUID customerId, CreateOrderRequest request) {
        BigDecimal distanceKm = pricingService.calculateDistance(
                request.getPickupLat().doubleValue(), request.getPickupLng().doubleValue(),
                request.getDropoffLat().doubleValue(), request.getDropoffLng().doubleValue());

        BigDecimal totalPrice = pricingService.calculatePrice(distanceKm);

        Order order = Order.builder()
                .customerId(customerId)
                .pickupAddress(request.getPickupAddress())
                .pickupLat(request.getPickupLat())
                .pickupLng(request.getPickupLng())
                .dropoffAddress(request.getDropoffAddress())
                .dropoffLat(request.getDropoffLat())
                .dropoffLng(request.getDropoffLng())
                .note(request.getNote())
                .distanceKm(distanceKm)
                .totalPrice(totalPrice)
                .status(OrderStatus.PENDING)
                .build();

        if (request.getItems() != null) {
            List<OrderItem> items = request.getItems().stream()
                    .map(i -> OrderItem.builder()
                            .order(order)
                            .name(i.getName())
                            .description(i.getDescription())
                            .quantity(i.getQuantity())
                            .weightKg(i.getWeightKg())
                            .imageUrl(i.getImageUrl())
                            .build())
                    .collect(Collectors.toList());
            order.getItems().addAll(items);
        }

        Order saved = orderRepository.save(order);
        eventPublisher.publishOrderCreated(saved);
        log.info("Order created: id={}, customer={}, price={}", saved.getId(), customerId, totalPrice);
        return toResponse(saved);
    }

    @Transactional
    public OrderResponse acceptOrder(UUID driverId, UUID orderId) {
        Order order = orderRepository.findById(orderId)
                .orElseThrow(() -> new OrderNotFoundException(orderId));

        if (order.getStatus() != OrderStatus.PENDING) {
            throw new OrderStatusException("Order is not available for acceptance. Current status: " + order.getStatus());
        }

        if (orderRepository.countActiveOrdersByDriver(driverId) > 0) {
            throw new OrderStatusException("Driver already has an active order");
        }

        order.setDriverId(driverId);
        order.setStatus(OrderStatus.ACCEPTED);
        Order saved = orderRepository.save(order);
        eventPublisher.publishOrderAccepted(saved);
        log.info("Order {} accepted by driver {}", orderId, driverId);
        return toResponse(saved);
    }

    @Transactional
    public OrderResponse updateStatus(UUID driverId, UUID orderId, OrderStatus newStatus) {
        Order order = orderRepository.findByIdAndDriverId(orderId, driverId)
                .orElseThrow(() -> new OrderNotFoundException(orderId));

        validateStatusTransition(order.getStatus(), newStatus);
        order.setStatus(newStatus);
        Order saved = orderRepository.save(order);

        if (newStatus == OrderStatus.DELIVERED) {
            eventPublisher.publishOrderDelivered(saved);
        }

        log.info("Order {} status updated to {} by driver {}", orderId, newStatus, driverId);
        return toResponse(saved);
    }

    @Transactional
    public OrderResponse cancelOrder(UUID customerId, UUID orderId, String reason) {
        Order order = orderRepository.findByIdAndCustomerId(orderId, customerId)
                .orElseThrow(() -> new OrderNotFoundException(orderId));

        if (order.getStatus() == OrderStatus.DELIVERED || order.getStatus() == OrderStatus.CANCELLED) {
            throw new OrderStatusException("Cannot cancel order with status: " + order.getStatus());
        }

        order.setStatus(OrderStatus.CANCELLED);
        order.setCancelledReason(reason);
        Order saved = orderRepository.save(order);
        eventPublisher.publishOrderCancelled(saved);
        log.info("Order {} cancelled by customer {}", orderId, customerId);
        return toResponse(saved);
    }

    public OrderResponse getOrder(UUID orderId) {
        return toResponse(orderRepository.findById(orderId)
                .orElseThrow(() -> new OrderNotFoundException(orderId)));
    }

    public Page<OrderResponse> getCustomerOrders(UUID customerId, Pageable pageable) {
        return orderRepository.findByCustomerIdOrderByCreatedAtDesc(customerId, pageable)
                .map(this::toResponse);
    }

    public Page<OrderResponse> getDriverOrders(UUID driverId, Pageable pageable) {
        return orderRepository.findByDriverIdOrderByCreatedAtDesc(driverId, pageable)
                .map(this::toResponse);
    }

    public List<OrderResponse> getPendingOrders() {
        return orderRepository.findPendingOrders().stream()
                .map(this::toResponse)
                .collect(Collectors.toList());
    }

    private void validateStatusTransition(OrderStatus current, OrderStatus next) {
        boolean valid = switch (current) {
            case ACCEPTED -> next == OrderStatus.PICKUP;
            case PICKUP -> next == OrderStatus.IN_TRANSIT;
            case IN_TRANSIT -> next == OrderStatus.DELIVERED;
            default -> false;
        };
        if (!valid) {
            throw new OrderStatusException("Invalid status transition: " + current + " → " + next);
        }
    }

    private OrderResponse toResponse(Order order) {
        List<OrderResponse.OrderItemResponse> items = order.getItems().stream()
                .map(i -> OrderResponse.OrderItemResponse.builder()
                        .id(i.getId())
                        .name(i.getName())
                        .description(i.getDescription())
                        .quantity(i.getQuantity())
                        .weightKg(i.getWeightKg())
                        .imageUrl(i.getImageUrl())
                        .build())
                .collect(Collectors.toList());

        return OrderResponse.builder()
                .id(order.getId())
                .customerId(order.getCustomerId())
                .driverId(order.getDriverId())
                .pickupAddress(order.getPickupAddress())
                .pickupLat(order.getPickupLat())
                .pickupLng(order.getPickupLng())
                .dropoffAddress(order.getDropoffAddress())
                .dropoffLat(order.getDropoffLat())
                .dropoffLng(order.getDropoffLng())
                .status(order.getStatus().name())
                .totalPrice(order.getTotalPrice())
                .distanceKm(order.getDistanceKm())
                .note(order.getNote())
                .cancelledReason(order.getCancelledReason())
                .items(items)
                .createdAt(order.getCreatedAt())
                .updatedAt(order.getUpdatedAt())
                .build();
    }
}
