package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type PushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type NotificationService struct {
	projectID   string
	credentials string
}

func NewNotificationService(projectID, credentials string) *NotificationService {
	return &NotificationService{projectID: projectID, credentials: credentials}
}

func (s *NotificationService) SendPush(fcmToken, title, body string) error {
	if fcmToken == "" {
		return nil
	}

	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"token": fcmToken,
			"notification": map[string]string{
				"title": title,
				"body":  body,
			},
		},
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", s.projectID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	accessToken := os.Getenv("FIREBASE_ACCESS_TOKEN")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("FCM returned status %d for token %s", resp.StatusCode, fcmToken)
	}
	return nil
}

func (s *NotificationService) HandleOrderCreated(orderID, customerFcmToken string) {
	s.SendPush(customerFcmToken, "Order Placed", "Your order has been placed successfully!")
}

func (s *NotificationService) HandleOrderAccepted(orderID, customerFcmToken, driverName string) {
	s.SendPush(customerFcmToken, "Driver Found", fmt.Sprintf("%s is on the way to pick up your order.", driverName))
}

func (s *NotificationService) HandleOrderCompleted(orderID, customerFcmToken string) {
	s.SendPush(customerFcmToken, "Order Delivered", "Your order has been delivered. Thank you!")
}

func (s *NotificationService) HandleOrderCancelled(orderID, customerFcmToken string) {
	s.SendPush(customerFcmToken, "Order Cancelled", "Your order has been cancelled.")
}

func (s *NotificationService) HandlePaymentCompleted(orderID, customerFcmToken string, amount float64) {
	s.SendPush(customerFcmToken, "Payment Successful", fmt.Sprintf("Payment of %.2f THB completed.", amount))
}
