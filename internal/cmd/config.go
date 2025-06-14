package cmd

import (
	"bytes"

	"github.com/dhth/act3/internal/types"
	"gopkg.in/yaml.v3"
)

type config struct {
	Workflows []types.Workflow `yaml:"workflows"`
}

func getWorkflowFromConfig(configBytes []byte) ([]types.Workflow, error) {
	config := config{}
	err := yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return config.Workflows, nil
}

func (c config) MarshalToYAML() ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(&c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
