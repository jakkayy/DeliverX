package main

import (
	"log"

	"github.com/deliverx/payment-service/config"
	"github.com/deliverx/payment-service/internal/handler"
	"github.com/deliverx/payment-service/internal/model"
	"github.com/deliverx/payment-service/internal/repository"
	"github.com/deliverx/payment-service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&model.Payment{})

	paymentRepo := repository.NewPaymentRepository(db)
	paymentSvc := service.NewPaymentService(paymentRepo, cfg.KafkaBrokers)
	paymentHandler := handler.NewPaymentHandler(paymentSvc, cfg.JWTSecret)

	r := gin.Default()
	paymentHandler.RegisterRoutes(r)

	log.Printf("payment-service running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
