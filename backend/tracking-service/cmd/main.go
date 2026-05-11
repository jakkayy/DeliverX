package main

import (
	"log"

	"github.com/deliverx/tracking-service/config"
	"github.com/deliverx/tracking-service/internal/handler"
	"github.com/deliverx/tracking-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})

	trackingSvc := service.NewTrackingService(rdb)
	trackingHandler := handler.NewTrackingHandler(trackingSvc, cfg.JWTSecret)

	r := gin.Default()
	trackingHandler.RegisterRoutes(r)

	log.Printf("tracking-service running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
