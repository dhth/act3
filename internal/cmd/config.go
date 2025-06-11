package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dhth/act3/internal/types"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Workflows []types.Workflow `yaml:"workflows"`
}

func expandTilde(path string, homeDir string) string {
	pathWithoutTilde, found := strings.CutPrefix(path, "~/")
	if !found {
		return path
	}
	return filepath.Join(homeDir, pathWithoutTilde)
}

func ReadConfig(configFilePath, userHomeDir string) ([]types.Workflow, error) {
	localFile, err := os.ReadFile(expandTilde(configFilePath, userHomeDir))
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(localFile, &config)
	if err != nil {
		return nil, err
	}

	return config.Workflows, nil
}
