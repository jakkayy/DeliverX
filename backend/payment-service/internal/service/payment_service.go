package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/deliverx/payment-service/internal/model"
	"github.com/deliverx/payment-service/internal/repository"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type CreatePaymentRequest struct {
	OrderID    uuid.UUID          `json:"order_id" binding:"required"`
	Amount     float64            `json:"amount" binding:"required"`
	Method     model.PaymentMethod `json:"method" binding:"required"`
}

type PaymentService struct {
	repo        *repository.PaymentRepository
	kafkaWriter *kafka.Writer
}

func NewPaymentService(repo *repository.PaymentRepository, kafkaBrokers string) *PaymentService {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(kafkaBrokers, ",")...),
		Balancer: &kafka.LeastBytes{},
	}
	return &PaymentService{repo: repo, kafkaWriter: writer}
}

func (s *PaymentService) CreatePayment(customerID uuid.UUID, req CreatePaymentRequest) (*model.Payment, error) {
	existing, _ := s.repo.FindByOrderID(req.OrderID)
	if existing != nil {
		return nil, errors.New("payment already exists for this order")
	}

	payment := &model.Payment{
		OrderID:    req.OrderID,
		CustomerID: customerID,
		Amount:     req.Amount,
		Method:     req.Method,
		Status:     model.PaymentStatusPending,
	}

	if req.Method == model.PaymentMethodCash {
		payment.Status = model.PaymentStatusCompleted
	}

	if err := s.repo.Save(payment); err != nil {
		return nil, err
	}

	if payment.Status == model.PaymentStatusCompleted {
		s.publishEvent("payment.completed", payment)
	}

	return payment, nil
}

func (s *PaymentService) CompletePayment(orderID uuid.UUID, stripePaymentID string) (*model.Payment, error) {
	payment, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	payment.Status = model.PaymentStatusCompleted
	payment.StripePaymentID = stripePaymentID
	if err := s.repo.Update(payment); err != nil {
		return nil, err
	}
	s.publishEvent("payment.completed", payment)
	return payment, nil
}

func (s *PaymentService) FailPayment(orderID uuid.UUID) (*model.Payment, error) {
	payment, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	payment.Status = model.PaymentStatusFailed
	if err := s.repo.Update(payment); err != nil {
		return nil, err
	}
	s.publishEvent("payment.failed", payment)
	return payment, nil
}

func (s *PaymentService) GetByOrderID(orderID uuid.UUID) (*model.Payment, error) {
	return s.repo.FindByOrderID(orderID)
}

func (s *PaymentService) GetMyPayments(customerID uuid.UUID) ([]model.Payment, error) {
	return s.repo.FindByCustomerID(customerID)
}

func (s *PaymentService) publishEvent(topic string, payment *model.Payment) {
	payload, _ := json.Marshal(payment)
	s.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   []byte(payment.ID.String()),
		Value: payload,
	})
}
