package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	ApiKey      string
	AdminApiKey string
	Port        string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	environment := os.Getenv("ENV")

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT cannot be empty")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY cannot be empty")
	}

	adminApiKey := os.Getenv("ADMIN_API_KEY")

	return &Config{
		Env:         environment,
		ApiKey:      apiKey,
		AdminApiKey: adminApiKey,
		Port:        port,
	}, nil
}
