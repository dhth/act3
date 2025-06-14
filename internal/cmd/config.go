package cmd

import (
	"github.com/dhth/act3/internal/types"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Workflows []types.Workflow `yaml:"workflows"`
}

func ReadConfig(configBytes []byte) ([]types.Workflow, error) {
	config := Config{}
	err := yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return config.Workflows, nil
}
