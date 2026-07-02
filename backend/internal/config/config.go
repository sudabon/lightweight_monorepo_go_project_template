package config

import "os"

// Config はアプリケーション設定を表します。
type Config struct {
	AppEnv           string
	AppPort          string
	DatabaseURL      string
	CORSAllowOrigins string
}

// Load は環境変数からアプリケーション設定を読み込みます。
func Load() Config {
	return Config{
		AppEnv:           getEnv("APP_ENV", "local"),
		AppPort:          getEnv("APP_PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
