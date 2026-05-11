package handler

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/deliverx/tracking-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsClient struct {
	conn     *websocket.Conn
	orderID  string
}

type TrackingHandler struct {
	svc       *service.TrackingService
	jwtSecret []byte
	clients   map[string][]*wsClient
	mu        sync.RWMutex
}

func NewTrackingHandler(svc *service.TrackingService, jwtSecret string) *TrackingHandler {
	key, _ := base64.StdEncoding.DecodeString(jwtSecret)
	return &TrackingHandler{
		svc:       svc,
		jwtSecret: key,
		clients:   make(map[string][]*wsClient),
	}
}

func (h *TrackingHandler) UpdateLocation(c *gin.Context) {
	driverID := h.userIDFromToken(c)
	if driverID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	var body struct {
		Lat float64 `json:"lat" binding:"required"`
		Lng float64 `json:"lng" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateDriverLocation(driverID, body.Lat, body.Lng); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.broadcast(driverID, body.Lat, body.Lng)
	c.Status(http.StatusNoContent)
}

func (h *TrackingHandler) GetDriverLocation(c *gin.Context) {
	driverID := c.Param("driverId")
	loc, err := h.svc.GetDriverLocation(driverID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}
	c.JSON(http.StatusOK, loc)
}

func (h *TrackingHandler) GetNearbyDrivers(c *gin.Context) {
	var q struct {
		Lat    float64 `form:"lat" binding:"required"`
		Lng    float64 `form:"lng" binding:"required"`
		Radius float64 `form:"radius"`
	}
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if q.Radius == 0 {
		q.Radius = 5
	}
	drivers, err := h.svc.GetNearbyDrivers(q.Lat, q.Lng, q.Radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"drivers": drivers})
}

func (h *TrackingHandler) WatchOrder(c *gin.Context) {
	orderID := c.Param("orderId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	client := &wsClient{conn: conn, orderID: orderID}
	h.mu.Lock()
	h.clients[orderID] = append(h.clients[orderID], client)
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		clients := h.clients[orderID]
		for i, cl := range clients {
			if cl == client {
				h.clients[orderID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		h.mu.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *TrackingHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "tracking-service"})
}

func (h *TrackingHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/tracking")
	{
		v1.GET("/health", h.Health)
		v1.POST("/location", h.UpdateLocation)
		v1.GET("/driver/:driverId", h.GetDriverLocation)
		v1.GET("/nearby", h.GetNearbyDrivers)
		v1.GET("/ws/order/:orderId", h.WatchOrder)
	}
}

func (h *TrackingHandler) broadcast(driverID string, lat, lng float64) {
	payload, _ := json.Marshal(map[string]interface{}{
		"driver_id": driverID,
		"lat":       lat,
		"lng":       lng,
	})
	h.mu.RLock()
	defer h.mu.RUnlock()

	driverKey, _ := h.svc.GetOrderDriver(driverID)
	for _, client := range h.clients[driverKey] {
		client.conn.WriteMessage(websocket.TextMessage, payload)
	}
}

func (h *TrackingHandler) userIDFromToken(c *gin.Context) string {
	tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if tokenStr == "" {
		return ""
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	sub, _ := claims["sub"].(string)
	return sub
}
