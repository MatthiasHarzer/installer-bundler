package core

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"installer-bundler/util/fsutil"
)

const configFile = "config/config.json"

type configItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type config struct {
	Items []configItem `json:"items"`
}

func (b *Bundler) BuildProject(destinationFile string) error {
	buildDir, cleanup, err := fsutil.CreateTempDirectory()
	if err != nil {
		return err
	}
	defer cleanup()

	err = fsutil.CopyFS(buildDir, b.runtimeProjectFS)
	if err != nil {
		return fmt.Errorf("failed to copy runtime project files: %w", err)
	}

	cfg := config{}
	for _, item := range b.items {
		cfg.Items = append(cfg.Items, configItem{
			Name: item.Title,
			URL:  item.Link,
		})
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	configFilePath := fmt.Sprintf("%s/%s", buildDir, configFile)

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	destinationDir := path.Dir(destinationFile)
	err = os.MkdirAll(destinationDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	cmd := exec.Command("make", "build")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	outputFile := fmt.Sprintf("%s/build/installer-runtime.exe", buildDir)

	err = fsutil.MoveFile(outputFile, destinationFile)
	if err != nil {
		return fmt.Errorf("failed to move built file to destination: %w", err)
	}

	return nil
}
