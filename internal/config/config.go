package config

import (
	"log"
	"os"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

type CurrencyConfig struct {
	RealTime string
	History  string
}

func LoadConfig() Config {
	required := []string{"DB_USER", "DB_PASS", "DB_HOST", "DB_PORT", "DB_NAME"}
	for _, v := range required {
		if os.Getenv(v) == "" {
			log.Fatal("Value not defined: ", v)
		}
	}

	return Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
	}
}

func GetApiKey() string {
	key := os.Getenv("EXCHANGE_KEY")
	if key == "" {
		log.Fatal("EXCHANGE_KEY not defined.")
	}
	return key
}

func GetExchangeList() CurrencyConfig {
	return CurrencyConfig{
		RealTime: os.Getenv("EXCHANGE_RT"),
		History:  os.Getenv("EXCHANGE_HIST"),
	}
}
