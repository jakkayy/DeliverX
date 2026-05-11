package config

import "os"

type Config struct {
	Port              string
	DBUrl             string
	JWTSecret         string
	KafkaBrokers      string
	StripeSecretKey   string
	StripeWebhookSecret string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8085"),
		DBUrl:               getEnv("DATABASE_URL", "host=localhost user=grab_user password=grab_secret_password dbname=grab_db port=5432 sslmode=disable"),
		JWTSecret:           getEnv("JWT_SECRET", "bXlfc3VwZXJfc2VjcmV0X2tleV9mb3JfZ3JhYl9kZWxpdmVyeV9hcHBfMjAyNA=="),
		KafkaBrokers:        getEnv("KAFKA_BROKERS", "localhost:29092"),
		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", "sk_test_placeholder"),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", "whsec_placeholder"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
