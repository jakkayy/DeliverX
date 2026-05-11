package repository

import (
	"github.com/deliverx/order-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Save(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id uuid.UUID) (*model.Order, error) {
	var order model.Order
	if err := r.db.First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) FindByCustomerID(customerID uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindByDriverID(driverID uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Where("driver_id = ?", driverID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindPendingOrders() ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Where("status = ?", model.StatusPending).Order("created_at ASC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
