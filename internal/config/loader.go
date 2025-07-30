package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadPacketConfig читает и парсит файл packet.json.
func LoadPacketConfig(filePath string) (*PacketConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %s: %w", filePath, err)
	}

	var cfg PacketConfig

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file %s: %w", filePath, err)
	}

	return &cfg, nil
}

func LoadUpdateConfig(filePath string) (*UpdateConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read update config file %s: %w", filePath, err)
	}

	var cfg UpdateConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse update config file %s: %w", filePath, err)
	}

	return &cfg, nil
}
