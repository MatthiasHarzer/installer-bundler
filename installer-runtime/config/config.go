package config

import (
	_ "embed"
	"encoding/json"
	"os"
)

//go:embed config.json
var configData []byte

type Mode string

const (
	ModeURL      Mode = "url"
	ModeEmbedded Mode = "embedded"
)

type Item struct {
	Name string  `json:"name"`
	URL  *string `json:"url"`
	File *string `json:"file"`
}

type Config struct {
	Items []*Item `json:"items"`
	Mode  Mode    `json:"mode"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := json.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func GetConfigFromFile(path string) (*Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
