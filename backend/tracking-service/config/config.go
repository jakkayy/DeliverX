package config

import "os"

type Config struct {
	Port          string
	RedisAddr     string
	RedisPassword string
	JWTSecret     string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8084"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", "redis_secret_password"),
		JWTSecret:     getEnv("JWT_SECRET", "bXlfc3VwZXJfc2VjcmV0X2tleV9mb3JfZ3JhYl9kZWxpdmVyeV9hcHBfMjAyNA=="),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
