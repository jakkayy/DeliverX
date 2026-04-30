package com.grab.order.service;

import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;

@Service
public class PricingService {

    private static final double EARTH_RADIUS_KM = 6371.0;

    private static final BigDecimal BASE_FARE = BigDecimal.valueOf(40);
    private static final BigDecimal RATE_PER_KM = BigDecimal.valueOf(12);
    private static final BigDecimal MIN_FARE = BigDecimal.valueOf(50);

    public BigDecimal calculateDistance(double lat1, double lng1, double lat2, double lng2) {
        double dLat = Math.toRadians(lat2 - lat1);
        double dLng = Math.toRadians(lng2 - lng1);
        double a = Math.sin(dLat / 2) * Math.sin(dLat / 2)
                + Math.cos(Math.toRadians(lat1)) * Math.cos(Math.toRadians(lat2))
                * Math.sin(dLng / 2) * Math.sin(dLng / 2);
        double c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
        double distanceKm = EARTH_RADIUS_KM * c;
        return BigDecimal.valueOf(distanceKm).setScale(2, RoundingMode.HALF_UP);
    }

    public BigDecimal calculatePrice(BigDecimal distanceKm) {
        BigDecimal price = BASE_FARE.add(RATE_PER_KM.multiply(distanceKm));
        return price.max(MIN_FARE).setScale(2, RoundingMode.HALF_UP);
    }
}
