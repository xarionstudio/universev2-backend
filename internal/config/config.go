package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName              string
	AppEnv               string
	AppPort              string
	AppURL               string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	DBMaxPoolSize        int
	DBMinPoolSize        int
	JWTSecret            string
	JWTExpiration        time.Duration
	JWTRefreshExpiration time.Duration
	FingerprintPort      string
	FingerprintBaud      int
	FingerprintEnabled   bool
	CORSAllowedOrigins   string
	LogLevel             string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppName:            getEnv("APP_NAME", "universev2-backend"),
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		AppURL:             getEnv("APP_URL", "http://localhost:8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "universev2"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		DBMaxPoolSize:      getEnvInt("DB_MAX_POOL_SIZE", 25),
		DBMinPoolSize:      getEnvInt("DB_MIN_POOL_SIZE", 5),
		JWTSecret:          getEnv("JWT_SECRET", "xxcXfZH1kaSiAwiO7wluFSjfYp8fj5+AbMvdssfEfT0xRtFXrpsA7H/I5BGnpMf1"),
		FingerprintPort:    getEnv("FINGERPRINT_PORT", "/dev/ttyUSB0"),
		FingerprintBaud:    getEnvInt("FINGERPRINT_BAUD", 115200),
		FingerprintEnabled: getEnvBool("FINGERPRINT_ENABLED", false),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		LogLevel:           getEnv("LOG_LEVEL", "debug"),
	}

	jwtExp, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		jwtExp = 24 * time.Hour
	}
	cfg.JWTExpiration = jwtExp

	jwtRefreshExp, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRATION", "168h"))
	if err != nil {
		jwtRefreshExp = 168 * time.Hour
	}
	cfg.JWTRefreshExpiration = jwtRefreshExp

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}
