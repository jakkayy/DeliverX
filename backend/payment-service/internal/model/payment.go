package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"

	PaymentMethodCard   PaymentMethod = "CARD"
	PaymentMethodWallet PaymentMethod = "WALLET"
	PaymentMethodCash   PaymentMethod = "CASH"
)

type Payment struct {
	ID              uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID         uuid.UUID     `gorm:"type:uuid;not null;column:order_id;uniqueIndex" json:"order_id"`
	CustomerID      uuid.UUID     `gorm:"type:uuid;not null;column:customer_id" json:"customer_id"`
	Amount          float64       `gorm:"not null" json:"amount"`
	Status          PaymentStatus `gorm:"size:20;not null;default:'PENDING'" json:"status"`
	Method          PaymentMethod `gorm:"size:20;not null" json:"method"`
	StripePaymentID string        `gorm:"column:stripe_payment_id;size:200" json:"stripe_payment_id,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

func (Payment) TableName() string { return "payments" }
