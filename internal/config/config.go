package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUser string
	DatabasePass string
	DatabaseHost string
	DatabasePort string
	DatabaseName string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load .env file. Verify if variables are defined.")
	}

	return Config{
		DatabaseUser: os.Getenv("DB_USER"),
		DatabasePass: os.Getenv("DB_PASS"),
		DatabaseHost: os.Getenv("DB_HOST"),
		DatabasePort: os.Getenv("DB_PORT"),
		DatabaseName: os.Getenv("DB_NAME"),
	}
}
