package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	OllamaURL   string
	EmbedModel  string
	LLMModel    string
	Port        string
}

func Load() Config {
	return Config{
		DatabaseURL: env("DATABASE_URL", "postgres://rag:rag@localhost:5432/rag?sslmode=disable"),
		OllamaURL:   env("OLLAMA_URL", "http://localhost:11435"),
		EmbedModel:  env("EMBED_MODEL", "bge-m3"),
		LLMModel:    env("LLM_MODEL", "llama3.2"),
		Port:        env("PORT", "8080"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
