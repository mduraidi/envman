package internal

import (
	"encoding/json"
	"os"
)

type EnvConfig struct {
	Dotnet string `json:"dotnet"`
	Python string `json:"python"`
	Nodejs string `json:"nodejs"`
	Golang string `json:"golang"`
	Java   string `json:"java"`
}

func LoadEnvConfig(path string) (*EnvConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg EnvConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveEnvConfig(path string, cfg *EnvConfig) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
