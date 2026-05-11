package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deliverx/notification-service/config"
	"github.com/deliverx/notification-service/internal/consumer"
	"github.com/deliverx/notification-service/internal/handler"
	"github.com/deliverx/notification-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	notifSvc := service.NewNotificationService(cfg.FirebaseProjectID, cfg.FirebaseCredentials)

	topics := []string{"order.created", "order.accepted", "order.completed", "order.cancelled", "payment.completed"}
	eventConsumer := consumer.NewEventConsumer(notifSvc, cfg.KafkaBrokers, "notification-service-group", topics)

	ctx, cancel := context.WithCancel(context.Background())
	go eventConsumer.Start(ctx)

	r := gin.Default()
	handler.RegisterRoutes(r)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("notification-service running on :%s", cfg.Port)
		if err := r.Run(":" + cfg.Port); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-quit
	cancel()
	log.Println("notification-service shutting down")
}
