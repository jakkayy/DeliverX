package config

import "os"

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8082"),
		DBUrl:     getEnv("DATABASE_URL", "host=localhost user=grab_user password=grab_secret_password dbname=grab_db port=5432 sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "bXlfc3VwZXJfc2VjcmV0X2tleV9mb3JfZ3JhYl9kZWxpdmVyeV9hcHBfMjAyNA=="),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
