package handler

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/deliverx/payment-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	svc       *service.PaymentService
	jwtSecret []byte
}

func NewPaymentHandler(svc *service.PaymentService, jwtSecret string) *PaymentHandler {
	key, _ := base64.StdEncoding.DecodeString(jwtSecret)
	return &PaymentHandler{svc: svc, jwtSecret: key}
}

func (h *PaymentHandler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return h.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["sub"])
		c.Next()
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	customerID, err := h.userIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	var req service.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payment, err := h.svc.CreatePayment(customerID, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) GetByOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	payment, err := h.svc.GetByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetMyPayments(c *gin.Context) {
	customerID, _ := h.userIDFromContext(c)
	payments, err := h.svc.GetMyPayments(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "payment-service"})
}

func (h *PaymentHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/payments")
	v1.GET("/health", h.Health)

	auth := v1.Use(h.authMiddleware())
	{
		auth.POST("", h.CreatePayment)
		auth.GET("", h.GetMyPayments)
		auth.GET("/order/:orderId", h.GetByOrder)
	}
}

func (h *PaymentHandler) userIDFromContext(c *gin.Context) (uuid.UUID, error) {
	raw, _ := c.Get("user_id")
	return uuid.Parse(raw.(string))
}
