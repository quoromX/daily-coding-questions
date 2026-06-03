package config

import (
	"os"
	"strings"
)

type Config struct {
	AppName            string
	AppEnv             string
	AppURL             string
	APIURL             string
	DatabaseURL        string
	JWTAccessSecret    string
	JWTRefreshSecret   string
	PasswordPepper     string
	Judge0BaseURL      string
	Judge0AuthToken    string
	Judge0WebhookToken string
	CORSOrigins        []string
	Port               string
	LogLevel           string
}

func Load() Config {
	return Config{
		AppName:            value("APP_NAME", "DaoForge"),
		AppEnv:             value("APP_ENV", "development"),
		AppURL:             value("APP_URL", "http://localhost:3000"),
		APIURL:             value("API_URL", "http://localhost:8080"),
		DatabaseURL:        value("DATABASE_URL", ""),
		JWTAccessSecret:    value("JWT_ACCESS_SECRET", "dev-access-secret-change-me"),
		JWTRefreshSecret:   value("JWT_REFRESH_SECRET", "dev-refresh-secret-change-me"),
		PasswordPepper:     value("PASSWORD_PEPPER", "dev-pepper-change-me"),
		Judge0BaseURL:      strings.TrimRight(value("JUDGE0_BASE_URL", ""), "/"),
		Judge0AuthToken:    value("JUDGE0_AUTH_TOKEN", ""),
		Judge0WebhookToken: value("JUDGE0_WEBHOOK_SECRET", ""),
		CORSOrigins:        split(value("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		Port:               value("PORT", "8080"),
		LogLevel:           value("LOG_LEVEL", "debug"),
	}
}

func value(key, fallback string) string {
	if got := os.Getenv(key); got != "" {
		return got
	}
	return fallback
}

func split(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
