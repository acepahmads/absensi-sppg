package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	JWTSecret  string
	AppPort    string
}

func LoadConfig() Config {
	// Load file .env
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Gagal load file .env: %v", err)
	}

	return Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		AppPort:    os.Getenv("APP_PORT"),
	}
}
