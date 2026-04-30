package com.grab.order.repository;

import com.grab.order.model.Order;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface OrderRepository extends JpaRepository<Order, UUID> {

    Page<Order> findByCustomerIdOrderByCreatedAtDesc(UUID customerId, Pageable pageable);

    Page<Order> findByDriverIdOrderByCreatedAtDesc(UUID driverId, Pageable pageable);

    @Query("SELECT o FROM Order o WHERE o.status = 'PENDING' ORDER BY o.createdAt ASC")
    List<Order> findPendingOrders();

    Optional<Order> findByIdAndCustomerId(UUID id, UUID customerId);

    Optional<Order> findByIdAndDriverId(UUID id, UUID driverId);

    @Query("SELECT COUNT(o) FROM Order o WHERE o.driverId = :driverId AND o.status NOT IN ('DELIVERED', 'CANCELLED')")
    long countActiveOrdersByDriver(UUID driverId);
}
