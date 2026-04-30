package com.grab.user.repository;

import com.grab.user.model.Driver;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface DriverRepository extends JpaRepository<Driver, UUID> {
    Optional<Driver> findByUserId(UUID userId);
    boolean existsByUserId(UUID userId);

    @Query("SELECT d FROM Driver d WHERE d.isAvailable = true")
    List<Driver> findAllAvailable();
}
