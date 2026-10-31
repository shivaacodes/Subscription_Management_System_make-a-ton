package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port               string
	Env                string
	DatabaseURL        string
	SessionSecret      string
	CORSAllowedOrigins []string
	CookieSecure       bool
	CookieDomain       string
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "development")
	dbURL := getEnv("DATABASE_URL", "postgres://subflow:subflow_secret@localhost:5432/subflow?sslmode=disable")
	sessionSecret := getEnv("SESSION_SECRET", "default_insecure_development_session_secret_32_chars_min!")

	corsRaw := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173")
	var origins []string
	for _, o := range strings.Split(corsRaw, ",") {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	cookieSecure := env == "production"
	if secureStr, exists := os.LookupEnv("COOKIE_SECURE"); exists {
		if b, err := strconv.ParseBool(secureStr); err == nil {
			cookieSecure = b
		}
	}

	return &Config{
		Port:               port,
		Env:                env,
		DatabaseURL:        dbURL,
		SessionSecret:      sessionSecret,
		CORSAllowedOrigins: origins,
		CookieSecure:       cookieSecure,
		CookieDomain:       os.Getenv("COOKIE_DOMAIN"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
