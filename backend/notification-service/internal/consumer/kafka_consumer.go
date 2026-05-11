package consumer

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/deliverx/notification-service/internal/service"
	"github.com/segmentio/kafka-go"
)

type EventConsumer struct {
	svc    *service.NotificationService
	reader *kafka.Reader
}

func NewEventConsumer(svc *service.NotificationService, brokers, groupID string, topics []string) *EventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokers, ","),
		GroupID: groupID,
		GroupTopics: topics,
	})
	return &EventConsumer{svc: svc, reader: reader}
}

func (c *EventConsumer) Start(ctx context.Context) {
	log.Println("notification consumer started")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("consumer error: %v", err)
			continue
		}
		c.handle(msg)
	}
}

func (c *EventConsumer) handle(msg kafka.Message) {
	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		log.Printf("failed to parse message: %v", err)
		return
	}

	orderID, _ := payload["id"].(string)
	fcmToken, _ := payload["fcm_token"].(string)

	switch msg.Topic {
	case "order.created":
		c.svc.HandleOrderCreated(orderID, fcmToken)
	case "order.accepted":
		c.svc.HandleOrderAccepted(orderID, fcmToken, "")
	case "order.completed":
		c.svc.HandleOrderCompleted(orderID, fcmToken)
	case "order.cancelled":
		c.svc.HandleOrderCancelled(orderID, fcmToken)
	case "payment.completed":
		amount, _ := payload["amount"].(float64)
		c.svc.HandlePaymentCompleted(orderID, fcmToken, amount)
	}
}
