package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Address              string `json:"address"`
	Port                 string `json:"port"`
	LogFile              string `json:"logFile"`
	LogLevel             string `json:"logLevel"`
	UseTLS               bool   `json:"useTLS"`
	UseMessageEncryption bool   `json:"useMessageEncryption"`
	InsecureSkipVerify   bool   `json:"insecureSkipVerify"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
