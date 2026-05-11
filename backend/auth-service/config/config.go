package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                string
	DBUrl               string
	RedisAddr           string
	RedisPassword       string
	JWTSecret           string
	JWTExpirationMs     int64
	RefreshExpirationMs int64
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8081"),
		DBUrl:               getEnv("DATABASE_URL", "host=localhost user=grab_user password=grab_secret_password dbname=grab_db port=5432 sslmode=disable"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:       getEnv("REDIS_PASSWORD", "redis_secret_password"),
		JWTSecret:           getEnv("JWT_SECRET", "bXlfc3VwZXJfc2VjcmV0X2tleV9mb3JfZ3JhYl9kZWxpdmVyeV9hcHBfMjAyNA=="),
		JWTExpirationMs:     getEnvInt64("JWT_EXPIRATION_MS", 900000),
		RefreshExpirationMs: getEnvInt64("JWT_REFRESH_EXPIRATION_MS", 604800000),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return fallback
}
