package core

import (
	"encoding/json"
	"os"
)

const configFile = "config/config.json"

type configItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type config struct {
	Items []configItem `json:"items"`
}

func (b *Bundler) GenerateProject() (string, error) {
	cfg := config{}
	for _, item := range b.items {
		cfg.Items = append(cfg.Items, configItem{
			Name: item.Title,
			URL:  item.Link,
		})
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}

	configFilePath := b.runtimeDir + "/" + configFile

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return "", err
	}

	return b.runtimeDir, nil
}
