package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadPacketConfig читает и парсит файл packet.json.
func LoadPacketConfig(filePath string) (*PacketConfig, error) {
	// Читаем содержимое файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %s: %w", filePath, err)
	}

	// Создаем переменную, куда будем загружать данные
	var cfg PacketConfig

	// Парсим JSON в нашу структуру
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file %s: %w", filePath, err)
	}

	return &cfg, nil
}
