package config

import "os"

type Config struct {
	Port                string
	JWTSecret           string
	AuthServiceURL      string
	UserServiceURL      string
	OrderServiceURL     string
	TrackingServiceURL  string
	PaymentServiceURL   string
	NotificationServiceURL string
}

func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "8080"),
		JWTSecret:              getEnv("JWT_SECRET", "bXlfc3VwZXJfc2VjcmV0X2tleV9mb3JfZ3JhYl9kZWxpdmVyeV9hcHBfMjAyNA=="),
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		UserServiceURL:         getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		OrderServiceURL:        getEnv("ORDER_SERVICE_URL", "http://localhost:8083"),
		TrackingServiceURL:     getEnv("TRACKING_SERVICE_URL", "http://localhost:8084"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "http://localhost:8085"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8086"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
