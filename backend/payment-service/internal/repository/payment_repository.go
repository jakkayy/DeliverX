package repository

import (
	"github.com/deliverx/payment-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Save(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) FindByOrderID(orderID uuid.UUID) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.Where("order_id = ?", orderID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepository) Update(payment *model.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) FindByCustomerID(customerID uuid.UUID) ([]model.Payment, error) {
	var payments []model.Payment
	if err := r.db.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
