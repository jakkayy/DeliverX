package main

import (
	"log"
	"net/http"

	"github.com/deliverx/api-gateway/config"
	"github.com/deliverx/api-gateway/internal/middleware"
	"github.com/deliverx/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	auth := middleware.AuthRequired(cfg.JWTSecret)

	r := gin.Default()
	r.Use(middleware.CORS())

	// auth-service — public (no JWT required)
	r.Any("/api/v1/auth/*path", proxy.To(cfg.AuthServiceURL))

	// protected services
	r.Any("/api/v1/users/*path", auth, proxy.To(cfg.UserServiceURL))
	r.Any("/api/v1/orders/*path", auth, proxy.To(cfg.OrderServiceURL))
	r.Any("/api/v1/tracking/*path", auth, proxy.To(cfg.TrackingServiceURL))
	r.Any("/api/v1/payments/*path", auth, proxy.To(cfg.PaymentServiceURL))
	r.Any("/api/v1/notifications/*path", auth, proxy.To(cfg.NotificationServiceURL))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "api-gateway"})
	})

	log.Printf("api-gateway running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
