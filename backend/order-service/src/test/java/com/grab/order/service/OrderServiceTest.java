package com.grab.order.service;

import com.grab.order.dto.CreateOrderRequest;
import com.grab.order.event.OrderEventPublisher;
import com.grab.order.exception.OrderNotFoundException;
import com.grab.order.exception.OrderStatusException;
import com.grab.order.model.Order;
import com.grab.order.model.Order.OrderStatus;
import com.grab.order.repository.OrderRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.math.BigDecimal;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class OrderServiceTest {

    @Mock private OrderRepository orderRepository;
    @Mock private OrderEventPublisher eventPublisher;
    @Mock private PricingService pricingService;
    @InjectMocks private OrderService orderService;

    private UUID customerId;
    private UUID driverId;
    private UUID orderId;
    private Order mockOrder;

    @BeforeEach
    void setUp() {
        customerId = UUID.randomUUID();
        driverId = UUID.randomUUID();
        orderId = UUID.randomUUID();

        mockOrder = Order.builder()
                .id(orderId)
                .customerId(customerId)
                .pickupAddress("123 Sukhumvit")
                .pickupLat(BigDecimal.valueOf(13.7563))
                .pickupLng(BigDecimal.valueOf(100.5018))
                .dropoffAddress("456 Silom")
                .dropoffLat(BigDecimal.valueOf(13.7300))
                .dropoffLng(BigDecimal.valueOf(100.5230))
                .status(OrderStatus.PENDING)
                .distanceKm(BigDecimal.valueOf(5.0))
                .totalPrice(BigDecimal.valueOf(100.0))
                .build();
    }

    @Test
    @DisplayName("Create order successfully")
    void createOrder_success() {
        CreateOrderRequest request = new CreateOrderRequest();
        request.setPickupAddress("123 Sukhumvit");
        request.setPickupLat(BigDecimal.valueOf(13.7563));
        request.setPickupLng(BigDecimal.valueOf(100.5018));
        request.setDropoffAddress("456 Silom");
        request.setDropoffLat(BigDecimal.valueOf(13.7300));
        request.setDropoffLng(BigDecimal.valueOf(100.5230));

        when(pricingService.calculateDistance(anyDouble(), anyDouble(), anyDouble(), anyDouble()))
                .thenReturn(BigDecimal.valueOf(5.0));
        when(pricingService.calculatePrice(any())).thenReturn(BigDecimal.valueOf(100.0));
        when(orderRepository.save(any(Order.class))).thenReturn(mockOrder);
        doNothing().when(eventPublisher).publishOrderCreated(any());

        var response = orderService.createOrder(customerId, request);

        assertThat(response.getStatus()).isEqualTo("PENDING");
        assertThat(response.getTotalPrice()).isEqualByComparingTo(BigDecimal.valueOf(100.0));
        verify(eventPublisher).publishOrderCreated(any());
    }

    @Test
    @DisplayName("Accept order successfully")
    void acceptOrder_success() {
        when(orderRepository.findById(orderId)).thenReturn(Optional.of(mockOrder));
        when(orderRepository.countActiveOrdersByDriver(driverId)).thenReturn(0L);
        when(orderRepository.save(any())).thenReturn(mockOrder);
        doNothing().when(eventPublisher).publishOrderAccepted(any());

        var response = orderService.acceptOrder(driverId, orderId);

        assertThat(response).isNotNull();
        verify(eventPublisher).publishOrderAccepted(any());
    }

    @Test
    @DisplayName("Accept order fails if driver has active order")
    void acceptOrder_driverBusy() {
        when(orderRepository.findById(orderId)).thenReturn(Optional.of(mockOrder));
        when(orderRepository.countActiveOrdersByDriver(driverId)).thenReturn(1L);

        assertThatThrownBy(() -> orderService.acceptOrder(driverId, orderId))
                .isInstanceOf(OrderStatusException.class)
                .hasMessageContaining("active order");
    }

    @Test
    @DisplayName("Cancel order successfully")
    void cancelOrder_success() {
        when(orderRepository.findByIdAndCustomerId(orderId, customerId))
                .thenReturn(Optional.of(mockOrder));
        when(orderRepository.save(any())).thenReturn(mockOrder);
        doNothing().when(eventPublisher).publishOrderCancelled(any());

        orderService.cancelOrder(customerId, orderId, "Changed mind");

        verify(eventPublisher).publishOrderCancelled(any());
    }

    @Test
    @DisplayName("Get order throws when not found")
    void getOrder_notFound() {
        when(orderRepository.findById(orderId)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> orderService.getOrder(orderId))
                .isInstanceOf(OrderNotFoundException.class);
    }

    @Test
    @DisplayName("Invalid status transition throws exception")
    void updateStatus_invalidTransition() {
        mockOrder.setStatus(OrderStatus.PENDING);
        mockOrder.setDriverId(driverId);

        when(orderRepository.findByIdAndDriverId(orderId, driverId))
                .thenReturn(Optional.of(mockOrder));

        assertThatThrownBy(() -> orderService.updateStatus(driverId, orderId, OrderStatus.DELIVERED))
                .isInstanceOf(OrderStatusException.class)
                .hasMessageContaining("Invalid status transition");
    }
}
