package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort         string
	DatabaseURL     string
	JWTSecret       string
	JWTTTLMin       int
	HospitalAAPIURL string
}

func Load() Config {
	return Config{
		AppPort:         getEnv("APP_PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://hospital_user:hospital_password@postgres:5432/hospital_db?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
		JWTTTLMin:       getEnvInt("JWT_TTL_MINUTES", 60),
		HospitalAAPIURL: getEnv("HOSPITAL_A_API_URL", "https://hospital-a.api.co.th"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
