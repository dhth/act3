package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/dhth/act3/internal/types"
	"gopkg.in/yaml.v3"
)

var (
	ErrConfigNotValid     = errors.New("config is not valid")
	errConfigNotValidYAML = errors.New("config is not valid YAML")
)

type config struct {
	Workflows []types.Workflow `yaml:"workflows"`
}

func getWorkflowsFromConfig(configBytes []byte) ([]types.Workflow, error) {
	var cfg config
	var zero []types.Workflow
	err := yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errConfigNotValidYAML, err.Error())
	}

	var errors []string
	for i, w := range cfg.Workflows {
		workflowErrors := w.Validate()
		if len(workflowErrors) > 0 {
			errors = append(errors, fmt.Sprintf("- workflow at index %d has errors: [%s]", i+1, strings.Join(workflowErrors, ", ")))
		}
	}

	if len(errors) > 0 {
		return zero, fmt.Errorf("%w:\n%s", ErrConfigNotValid, strings.Join(errors, "\n"))
	}

	return cfg.Workflows, nil
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
