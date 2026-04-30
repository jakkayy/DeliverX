package com.grab.user.service;

import com.grab.user.dto.DriverResponse;
import com.grab.user.dto.RegisterDriverRequest;
import com.grab.user.exception.ResourceNotFoundException;
import com.grab.user.model.Driver;
import com.grab.user.model.User;
import com.grab.user.repository.DriverRepository;
import com.grab.user.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class DriverService {

    private final DriverRepository driverRepository;
    private final UserRepository userRepository;

    @Transactional
    public DriverResponse registerDriver(UUID userId, RegisterDriverRequest request) {
        if (driverRepository.existsByUserId(userId)) {
            throw new IllegalArgumentException("Driver profile already exists for this user");
        }

        User user = userRepository.findById(userId)
                .orElseThrow(() -> new ResourceNotFoundException("User not found: " + userId));

        if (user.getRole() != User.Role.DRIVER) {
            throw new IllegalArgumentException("User role must be DRIVER to register as driver");
        }

        Driver driver = Driver.builder()
                .user(user)
                .vehicleType(Driver.VehicleType.valueOf(request.getVehicleType()))
                .licensePlate(request.getLicensePlate())
                .vehicleModel(request.getVehicleModel())
                .isAvailable(false)
                .build();

        driver = driverRepository.save(driver);
        log.info("Driver registered: userId={}", userId);
        return toResponse(driver);
    }

    public DriverResponse getDriverByUserId(UUID userId) {
        Driver driver = driverRepository.findByUserId(userId)
                .orElseThrow(() -> new ResourceNotFoundException("Driver not found for userId: " + userId));
        return toResponse(driver);
    }

    public DriverResponse getDriverById(UUID driverId) {
        Driver driver = driverRepository.findById(driverId)
                .orElseThrow(() -> new ResourceNotFoundException("Driver not found: " + driverId));
        return toResponse(driver);
    }

    @Transactional
    public DriverResponse updateAvailability(UUID userId, boolean isAvailable) {
        Driver driver = driverRepository.findByUserId(userId)
                .orElseThrow(() -> new ResourceNotFoundException("Driver not found for userId: " + userId));
        driver.setIsAvailable(isAvailable);
        log.info("Driver {} availability set to {}", driver.getId(), isAvailable);
        return toResponse(driverRepository.save(driver));
    }

    public List<DriverResponse> getAvailableDrivers() {
        return driverRepository.findAllAvailable().stream()
                .map(this::toResponse)
                .collect(Collectors.toList());
    }

    private DriverResponse toResponse(Driver driver) {
        return DriverResponse.builder()
                .id(driver.getId())
                .userId(driver.getUser().getId())
                .name(driver.getUser().getName())
                .phone(driver.getUser().getPhone())
                .profileImage(driver.getUser().getProfileImage())
                .vehicleType(driver.getVehicleType().name())
                .licensePlate(driver.getLicensePlate())
                .vehicleModel(driver.getVehicleModel())
                .rating(driver.getRating())
                .totalTrips(driver.getTotalTrips())
                .isAvailable(driver.getIsAvailable())
                .createdAt(driver.getCreatedAt())
                .build();
    }
}
