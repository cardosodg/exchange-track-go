package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseUser string
	DatabasePass string
	DatabaseHost string
	DatabasePort string
	DatabaseName string
}

func LoadConfig() Config {
	required := []string{"DB_USER", "DB_PASS", "DB_HOST", "DB_PORT", "DB_NAME", "API_KEY"}
	for _, v := range required {
		if os.Getenv(v) == "" {
			log.Printf("⚠️  Variável de ambiente %s não está definida", v)
		}
	}

	return Config{
		DatabaseUser: os.Getenv("DB_USER"),
		DatabasePass: os.Getenv("DB_PASS"),
		DatabaseHost: os.Getenv("DB_HOST"),
		DatabasePort: os.Getenv("DB_PORT"),
		DatabaseName: os.Getenv("DB_NAME"),
	}
}

func GetApiKey() string {
	key := os.Getenv("EXCHANGE_KEY")
	if key == "" {
		log.Fatal("EXCHANGE_KEY not defined.")
	}
	return key
}
