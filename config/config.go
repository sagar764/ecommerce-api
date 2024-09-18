package config

import (
	"ecommerce-api/internal/entities"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadConfig loads the configuration for the application based on the given appName.
// It uses environment variables and the "envconfig" package to populate the configuration struct.
// Parameters:
// - appName: The name of the application to load configuration for.
// Returns:
// - *entities.EnvConfig: A pointer to the populated configuration struct.
// - error: An error if there was an issue loading the configuration.
func LoadConfig(appName string) (*entities.EnvConfig, error) {

	var cfg entities.EnvConfig

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	err = envconfig.Process(appName, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
