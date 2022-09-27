package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env           string
	ApiKey        string
	AdminApiKey   string
	Port          string
	RedisHost     string
	RedisPassword string
}

func LoadConfig(env *string) (*Config, error) {

	err := godotenv.Load(*env)
	if err != nil {
		log.Println(".env file not found")
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

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		return nil, fmt.Errorf("REDIS_HOST cannot be empty")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		log.Println("REDIS_PASSWORD empty")
	}

	return &Config{
		Env:           environment,
		ApiKey:        apiKey,
		AdminApiKey:   adminApiKey,
		Port:          port,
		RedisHost:     redisHost,
		RedisPassword: redisPassword,
	}, nil
}
