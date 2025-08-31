package core

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	root "installer-bundler"

	cp "github.com/otiai10/copy"

	"installer-bundler/util/fsutil"
)

type Mode string

const (
	ModeURL      Mode = "url"
	ModeEmbedded Mode = "embedded"
)

const runtimeConfigFile = "config/config.json"
const runtimeFilesDir = "files"

type configItem struct {
	Name string  `json:"name"`
	URL  *string `json:"url,omitempty"`
	File *string `json:"file,omitempty"`
}

type config struct {
	Items []configItem `json:"items"`
	Mode  Mode         `json:"mode"`
}

func (b *Bundler) writeProjectFiles(buildDir string, cfg config) error {
	err := fsutil.CopyFS(buildDir, b.runtimeProjectFS)
	if err != nil {
		return fmt.Errorf("failed to copy runtime project files: %w", err)
	}

	buildRuntimeFilesDir := path.Join(buildDir, runtimeFilesDir)

	if cfg.Mode == ModeEmbedded {
		err = cp.Copy(b.fileCacheDir, buildRuntimeFilesDir)
		if err != nil {
			return fmt.Errorf("failed to copy files directory: %w", err)
		}
	} else {
		err = os.MkdirAll(buildRuntimeFilesDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create files directory: %w", err)
		}

		// This dummy file is needed for the directory to be embedded in the binary later
		dummyFile := path.Join(buildRuntimeFilesDir, "dummy")

		w, err := os.OpenFile(dummyFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	configFilePath := fmt.Sprintf("%s/%s", buildDir, runtimeConfigFile)

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bundler) build(cfg config, destinationFile string, copyItemFiles bool) error {
	buildDir, cleanup, err := fsutil.CreateTempDirectory()
	if err != nil {
		return err
	}
	defer cleanup()

	err = b.writeProjectFiles(buildDir, cfg)
	if err != nil {
		return err
	}

	destinationDir := path.Dir(destinationFile)
	err = os.MkdirAll(destinationDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	cmd := exec.Command("make", "build", fmt.Sprintf("BUILD_VESION=%s", root.Version))
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	var builtFile string
	if runtime.GOOS == "windows" {
		builtFile = fmt.Sprintf("%s/build/installer-runtime.exe", buildDir)
	} else {
		builtFile = fmt.Sprintf("%s/build/installer-runtime", buildDir)
	}

	err = fsutil.MoveFile(builtFile, destinationFile)
	if err != nil {
		return fmt.Errorf("failed to move built file to destination: %w", err)
	}

	return nil
}

func (b *Bundler) BuildProjectURL(destinationFile string) error {
	cfg := config{
		Mode: ModeURL,
	}
	for _, item := range b.items {
		cfg.Items = append(cfg.Items, configItem{
			Name: item.Title,
			URL:  &item.Link,
		})
	}

	return b.build(cfg, destinationFile, false)
}

func (b *Bundler) BuildProjectEmbedded(destinationFile string) error {
	cfg := config{
		Mode: ModeEmbedded,
	}
	for _, item := range b.items {
		isDownloaded, file := b.IsDownloaded(item)
		if !isDownloaded {
			var err error
			file, err = b.Download(item)
			if err != nil {
				return fmt.Errorf("failed to download item %s: %w", item.Title, err)
			}
		}

		cfg.Items = append(cfg.Items, configItem{
			Name: item.Title,
			File: &file,
		})
	}

	return b.build(cfg, destinationFile, true)
}
