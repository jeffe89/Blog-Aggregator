// Internal/config/config.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Structure for config file
type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Function to set current_user_name and write the config struct to a JSON file
func (cfg *Config) SetUser(userName string) error {

	// Update current user name
	cfg.CurrentUserName = userName

	// Write the updated config to disk
	return write(*cfg)
}

// Function to load config from disk
func Read() (Config, error) {

	// Get the path to config file
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	// Read the file contents
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	// Unmarshal the JSON into a Config struct
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Helper function to gather path to config file
func getConfigFilePath() (string, error) {

	// Get path to home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Join home directory with config file name and return
	return filepath.Join(homeDir, configFileName), nil
}

// Helper function to marshal Config struct to JSON and write to config file
func write(cfg Config) error {

	// Get path to the config file
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Marshal the config to JSON
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	// Write the JSON data to the config file
	return os.WriteFile(configPath, data, 0644)
}
