package main

import (
	"log"

	"github.com/deliverx/auth-service/config"
	"github.com/deliverx/auth-service/internal/handler"
	"github.com/deliverx/auth-service/internal/middleware"
	"github.com/deliverx/auth-service/internal/model"
	"github.com/deliverx/auth-service/internal/repository"
	"github.com/deliverx/auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&model.User{})

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})

	jwtUtil := middleware.NewJWTUtil(cfg.JWTSecret, cfg.JWTExpirationMs, cfg.RefreshExpirationMs)
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, jwtUtil, rdb)
	authHandler := handler.NewAuthHandler(authSvc)

	r := gin.Default()
	authHandler.RegisterRoutes(r)

	log.Printf("auth-service running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
