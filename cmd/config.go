package cmd

import (
	"os"
	"os/user"
	"strings"

	"github.com/dhth/act3/internal/types"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Workflows []types.Workflow `yaml:"workflows"`
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			os.Exit(1)
		}
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	return path
}

func ReadConfig(configFilePath string) ([]types.Workflow, error) {
	localFile, err := os.ReadFile(expandTilde(configFilePath))
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
