package config

import (
	"encoding/json"
	"fmt"
)

// PacketConfig представляет полную структуру файла packet.json
type PacketConfig struct {
	Name         string             `json:"name"`
	Version      string             `json:"ver"`
	Targets      []Target           `json:"targets"`
	Dependencies []PacketDependency `json:"packets"`
}

// PacketDependency представляет зависимость от другого пакета
type PacketDependency struct {
	Name    string `json:"name"`
	Version string `json:"ver"`
}

// Target представляет одну цель для архивации.
// Это наша "универсальная" структура, которая может хранить и строку, и объект.
type Target struct {
	Path    string
	Exclude string
}

// UnmarshalJSON - это специальный метод, который "учит" Go,
// как правильно парсить наше поле "targets" со смешанными типами.
func (t *Target) UnmarshalJSON(data []byte) error {
	// Сначала пробуем распарсить как простую строку
	var simplePath string
	if err := json.Unmarshal(data, &simplePath); err == nil {
		t.Path = simplePath
		return nil
	}

	// Если не получилось, пробуем распарсить как объект
	var complexTarget struct {
		Path    string `json:"path"`
		Exclude string `json:"exclude"`
	}
	if err := json.Unmarshal(data, &complexTarget); err == nil {
		t.Path = complexTarget.Path
		t.Exclude = complexTarget.Exclude
		return nil
	}

	return fmt.Errorf("target must be a string or a JSON object with 'path' and optional 'exclude' fields")
}
