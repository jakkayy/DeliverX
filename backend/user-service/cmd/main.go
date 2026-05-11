package main

import (
	"log"

	"github.com/deliverx/user-service/config"
	"github.com/deliverx/user-service/internal/handler"
	"github.com/deliverx/user-service/internal/model"
	"github.com/deliverx/user-service/internal/repository"
	"github.com/deliverx/user-service/internal/service"
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
	db.AutoMigrate(&model.User{})

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc, cfg.JWTSecret)

	r := gin.Default()
	userHandler.RegisterRoutes(r)

	log.Printf("user-service running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
