package config

import (
	_ "embed"
	"encoding/json"
)

//go:embed config.json
var configData []byte

type Item struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	Items []*Item `json:"items"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := json.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
