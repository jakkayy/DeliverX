package config

import "os"

type Config struct {
	Port                string
	KafkaBrokers        string
	FirebaseProjectID   string
	FirebaseCredentials string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8086"),
		KafkaBrokers:        getEnv("KAFKA_BROKERS", "localhost:29092"),
		FirebaseProjectID:   getEnv("FIREBASE_PROJECT_ID", "your-project-id"),
		FirebaseCredentials: getEnv("FIREBASE_CREDENTIALS_JSON", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
