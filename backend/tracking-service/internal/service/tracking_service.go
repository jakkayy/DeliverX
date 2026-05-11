package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	driverLocationPrefix = "driver:location:"
	driverGeoKey         = "drivers:geo"
	orderTrackingPrefix  = "order:tracking:"
)

type Location struct {
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TrackingService struct {
	redis *redis.Client
}

func NewTrackingService(rdb *redis.Client) *TrackingService {
	return &TrackingService{redis: rdb}
}

func (s *TrackingService) UpdateDriverLocation(driverID string, lat, lng float64) error {
	ctx := context.Background()
	loc := Location{Lat: lat, Lng: lng, UpdatedAt: time.Now()}
	payload, _ := json.Marshal(loc)

	s.redis.Set(ctx, driverLocationPrefix+driverID, payload, 10*time.Minute)
	s.redis.GeoAdd(ctx, driverGeoKey, &redis.GeoLocation{
		Name:      driverID,
		Longitude: lng,
		Latitude:  lat,
	})
	return nil
}

func (s *TrackingService) GetDriverLocation(driverID string) (*Location, error) {
	ctx := context.Background()
	val, err := s.redis.Get(ctx, driverLocationPrefix+driverID).Result()
	if err != nil {
		return nil, fmt.Errorf("driver location not found")
	}
	var loc Location
	json.Unmarshal([]byte(val), &loc)
	return &loc, nil
}

func (s *TrackingService) GetNearbyDrivers(lat, lng, radiusKm float64) ([]string, error) {
	ctx := context.Background()
	results, err := s.redis.GeoRadius(ctx, driverGeoKey, lng, lat, &redis.GeoRadiusQuery{
		Radius: radiusKm,
		Unit:   "km",
		Sort:   "ASC",
	}).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r.Name
	}
	return ids, nil
}

func (s *TrackingService) SetOrderTracking(orderID, driverID string) error {
	ctx := context.Background()
	return s.redis.Set(ctx, orderTrackingPrefix+orderID, driverID, 24*time.Hour).Err()
}

func (s *TrackingService) GetOrderDriver(orderID string) (string, error) {
	ctx := context.Background()
	return s.redis.Get(ctx, orderTrackingPrefix+orderID).Result()
}
