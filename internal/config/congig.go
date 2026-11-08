package config

import "os"

type Config struct {
	DBPath  string
	BaseURL string
}

func Load() Config {
	return Config{
		DBPath:  getenv("SHORTENER_DB", "shortener.db"),
		BaseURL: getenv("SHORTENER_BASE_URL", "http://localhost:8080"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
