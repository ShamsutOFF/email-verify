package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Email Email
}

type Email struct {
	Email    string
	Password string
	Address  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default configs")
	}
	return &Config{
		Email: Email{
			Email:    os.Getenv("EMAIL"),
			Password: os.Getenv("PASSWORD"),
			Address:  os.Getenv("ADDRESS"),
		},
	}
}
