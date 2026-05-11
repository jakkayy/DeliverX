package handler

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/deliverx/user-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc       *service.UserService
	jwtSecret []byte
}

func NewUserHandler(svc *service.UserService, jwtSecret string) *UserHandler {
	key, _ := base64.StdEncoding.DecodeString(jwtSecret)
	return &UserHandler{svc: svc, jwtSecret: key}
}

func (h *UserHandler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
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

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := h.userIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := h.userIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetDrivers(c *gin.Context) {
	drivers, err := h.svc.GetDrivers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, drivers)
}

func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "user-service"})
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/users")
	v1.GET("/health", h.Health)

	auth := v1.Use(h.authMiddleware())
	{
		auth.GET("/me", h.GetProfile)
		auth.PUT("/me", h.UpdateProfile)
		auth.GET("/drivers", h.GetDrivers)
	}
}

func (h *UserHandler) userIDFromContext(c *gin.Context) (uuid.UUID, error) {
	raw, _ := c.Get("user_id")
	return uuid.Parse(raw.(string))
}
