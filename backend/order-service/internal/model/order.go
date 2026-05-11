package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusAccepted   OrderStatus = "ACCEPTED"
	StatusPickedUp   OrderStatus = "PICKED_UP"
	StatusDelivering OrderStatus = "DELIVERING"
	StatusCompleted  OrderStatus = "COMPLETED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

type Order struct {
	ID               uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CustomerID       uuid.UUID   `gorm:"type:uuid;not null;column:customer_id" json:"customer_id"`
	DriverID         *uuid.UUID  `gorm:"type:uuid;column:driver_id" json:"driver_id,omitempty"`
	Status           OrderStatus `gorm:"size:20;not null;default:'PENDING'" json:"status"`
	PickupAddress    string      `gorm:"column:pickup_address;not null" json:"pickup_address"`
	PickupLat        float64     `gorm:"column:pickup_lat" json:"pickup_lat"`
	PickupLng        float64     `gorm:"column:pickup_lng" json:"pickup_lng"`
	DropoffAddress   string      `gorm:"column:dropoff_address;not null" json:"dropoff_address"`
	DropoffLat       float64     `gorm:"column:dropoff_lat" json:"dropoff_lat"`
	DropoffLng       float64     `gorm:"column:dropoff_lng" json:"dropoff_lng"`
	Note             string      `gorm:"size:500" json:"note,omitempty"`
	TotalPrice       float64     `gorm:"column:total_price" json:"total_price"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

func (Order) TableName() string { return "orders" }
