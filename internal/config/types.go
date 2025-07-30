package config

import (
	"encoding/json"
	"fmt"
)

type PacketConfig struct {
	Name         string             `json:"name"`
	Version      string             `json:"ver"`
	Targets      []Target           `json:"targets"`
	Dependencies []PacketDependency `json:"packets"`
}

type PacketDependency struct {
	Name    string `json:"name"`
	Version string `json:"ver"`
}

type UpdateConfig struct {
	Packages []PackageRequest `json:"packages"`
}

type PackageRequest struct {
	Name    string `json:"name"`
	Version string `json:"ver,omitempty"`
}

type Target struct {
	Path    string
	Exclude string
}

func (t *Target) UnmarshalJSON(data []byte) error {
	var simplePath string
	if err := json.Unmarshal(data, &simplePath); err == nil {
		t.Path = simplePath
		return nil
	}

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
