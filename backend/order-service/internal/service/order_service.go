package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/deliverx/order-service/internal/model"
	"github.com/deliverx/order-service/internal/repository"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type CreateOrderRequest struct {
	PickupAddress  string  `json:"pickup_address" binding:"required"`
	PickupLat      float64 `json:"pickup_lat"`
	PickupLng      float64 `json:"pickup_lng"`
	DropoffAddress string  `json:"dropoff_address" binding:"required"`
	DropoffLat     float64 `json:"dropoff_lat"`
	DropoffLng     float64 `json:"dropoff_lng"`
	Note           string  `json:"note"`
	TotalPrice     float64 `json:"total_price"`
}

type OrderService struct {
	repo        *repository.OrderRepository
	kafkaWriter *kafka.Writer
}

func NewOrderService(repo *repository.OrderRepository, kafkaBrokers string) *OrderService {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(kafkaBrokers, ",")...),
		Balancer: &kafka.LeastBytes{},
	}
	return &OrderService{repo: repo, kafkaWriter: writer}
}

func (s *OrderService) CreateOrder(customerID uuid.UUID, req CreateOrderRequest) (*model.Order, error) {
	order := &model.Order{
		CustomerID:     customerID,
		Status:         model.StatusPending,
		PickupAddress:  req.PickupAddress,
		PickupLat:      req.PickupLat,
		PickupLng:      req.PickupLng,
		DropoffAddress: req.DropoffAddress,
		DropoffLat:     req.DropoffLat,
		DropoffLng:     req.DropoffLng,
		Note:           req.Note,
		TotalPrice:     req.TotalPrice,
	}
	if err := s.repo.Save(order); err != nil {
		return nil, err
	}
	s.publishEvent("order.created", order)
	return order, nil
}

func (s *OrderService) GetOrder(orderID uuid.UUID) (*model.Order, error) {
	return s.repo.FindByID(orderID)
}

func (s *OrderService) GetMyOrders(customerID uuid.UUID) ([]model.Order, error) {
	return s.repo.FindByCustomerID(customerID)
}

func (s *OrderService) GetDriverOrders(driverID uuid.UUID) ([]model.Order, error) {
	return s.repo.FindByDriverID(driverID)
}

func (s *OrderService) GetPendingOrders() ([]model.Order, error) {
	return s.repo.FindPendingOrders()
}

func (s *OrderService) AcceptOrder(orderID, driverID uuid.UUID) (*model.Order, error) {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.Status != model.StatusPending {
		return nil, errors.New("order is not available")
	}
	order.Status = model.StatusAccepted
	order.DriverID = &driverID
	if err := s.repo.Update(order); err != nil {
		return nil, err
	}
	s.publishEvent("order.accepted", order)
	return order, nil
}

func (s *OrderService) UpdateStatus(orderID uuid.UUID, status model.OrderStatus) (*model.Order, error) {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	order.Status = status
	if err := s.repo.Update(order); err != nil {
		return nil, err
	}
	topic := "order." + strings.ToLower(string(status))
	s.publishEvent(topic, order)
	return order, nil
}

func (s *OrderService) publishEvent(topic string, order *model.Order) {
	payload, _ := json.Marshal(order)
	s.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   []byte(order.ID.String()),
		Value: payload,
	})
}
