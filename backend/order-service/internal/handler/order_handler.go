package handler

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/deliverx/order-service/internal/model"
	"github.com/deliverx/order-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	svc       *service.OrderService
	jwtSecret []byte
}

func NewOrderHandler(svc *service.OrderService, jwtSecret string) *OrderHandler {
	key, _ := base64.StdEncoding.DecodeString(jwtSecret)
	return &OrderHandler{svc: svc, jwtSecret: key}
}

func (h *OrderHandler) authMiddleware() gin.HandlerFunc {
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
		c.Set("role", claims["role"])
		c.Next()
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	customerID, err := h.userIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.CreateOrder(customerID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	order, err := h.svc.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	customerID, _ := h.userIDFromContext(c)
	orders, err := h.svc.GetMyOrders(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetPendingOrders(c *gin.Context) {
	orders, err := h.svc.GetPendingOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) AcceptOrder(c *gin.Context) {
	driverID, err := h.userIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	order, err := h.svc.AcceptOrder(orderID, driverID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	var body struct {
		Status model.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.UpdateStatus(orderID, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "order-service"})
}

func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/orders")
	v1.GET("/health", h.Health)

	auth := v1.Use(h.authMiddleware())
	{
		auth.POST("", h.CreateOrder)
		auth.GET("", h.GetMyOrders)
		auth.GET("/pending", h.GetPendingOrders)
		auth.GET("/:id", h.GetOrder)
		auth.POST("/:id/accept", h.AcceptOrder)
		auth.PUT("/:id/status", h.UpdateStatus)
	}
}

func (h *OrderHandler) userIDFromContext(c *gin.Context) (uuid.UUID, error) {
	raw, _ := c.Get("user_id")
	return uuid.Parse(raw.(string))
}
