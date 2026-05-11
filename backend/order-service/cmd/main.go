package main

import (
	"log"

	"github.com/deliverx/order-service/config"
	"github.com/deliverx/order-service/internal/handler"
	"github.com/deliverx/order-service/internal/model"
	"github.com/deliverx/order-service/internal/repository"
	"github.com/deliverx/order-service/internal/service"
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
	db.AutoMigrate(&model.Order{})

	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo, cfg.KafkaBrokers)
	orderHandler := handler.NewOrderHandler(orderSvc, cfg.JWTSecret)

	r := gin.Default()
	orderHandler.RegisterRoutes(r)

	log.Printf("order-service running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
