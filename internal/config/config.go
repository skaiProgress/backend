package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Port               string
	DatabaseURL        string
	AuthJWTSecret      string
	AuthJWTIssuer      string
	AuthJWTExpiry      time.Duration
	AppEnv             string
	CORSAllowedOrigins []string
	StoragePath        string
	PublicAPIURL       string
}

// Load reads configuration from environment variables.
// It attempts to load a .env file for local development; missing file is not an error.
func Load() (*Config, error) {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	corsOrigins := parseCSV(os.Getenv("CORS_ALLOWED_ORIGINS"))
	jwtExpiry := parseJWTExpiryHours(os.Getenv("AUTH_JWT_EXPIRY_HOURS"))

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}
	publicAPIURL := os.Getenv("PUBLIC_API_URL")
	if publicAPIURL == "" {
		publicAPIURL = "http://localhost:" + port
	}

	return &Config{
		Port:               port,
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		AuthJWTSecret:      os.Getenv("AUTH_JWT_SECRET"),
		AuthJWTIssuer:      os.Getenv("AUTH_JWT_ISSUER"),
		AuthJWTExpiry:      jwtExpiry,
		AppEnv:             appEnv,
		CORSAllowedOrigins: corsOrigins,
		StoragePath:        storagePath,
		PublicAPIURL:       publicAPIURL,
	}, nil
}

func parseJWTExpiryHours(raw string) time.Duration {
	if raw == "" {
		return 24 * time.Hour
	}
	hours, err := strconv.Atoi(raw)
	if err != nil || hours <= 0 {
		return 24 * time.Hour
	}
	return time.Duration(hours) * time.Hour
}

func parseCSV(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
