package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractRepoName(t *testing.T) {
	expectedRepo := "dhth/act3"
	testCases := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		// success
		{
			name:     "https url",
			input:    "https://github.com/dhth/act3.git",
			expected: expectedRepo,
		},
		{
			name:     "https url without .git suffix",
			input:    "https://github.com/dhth/act3",
			expected: expectedRepo,
		},
		{
			name:     "http url",
			input:    "http://github.com/dhth/act3.git",
			expected: expectedRepo,
		},
		{
			name:     "ssh url",
			input:    "ssh://git@github.com/dhth/act3.git",
			expected: expectedRepo,
		},
		{
			name:     "ssh url with port",
			input:    "ssh://git@github.com:443/dhth/act3.git",
			expected: expectedRepo,
		},
		{
			name:     "git protocol url",
			input:    "git://github.com/dhth/act3.git",
			expected: expectedRepo,
		},
		{
			name:     "git protocol url with port",
			input:    "git://github.com:9418/dhth/act3.git",
			expected: expectedRepo,
		},
		// failures
		{
			name:  "empty url",
			input: "",
			err:   errRemoteURLEmpty,
		},
		{
			name:  "invalid url",
			input: "invalid_url",
			err:   errInvalidURLFormat,
		},
		{
			name:  "url with no scheme",
			input: "github.com/dhth/act3.git",
			err:   errInvalidURLFormat,
		},
		{
			name:  "ssh url with subdirectory",
			input: "https://github.com/dhth/act3/subdir",
			err:   errInvalidURLFormat,
		},
		{
			name:  "https url with subdirectory",
			input: "https://github.com/dhth/act3/subdir",
			err:   errInvalidURLFormat,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractRepoName(tt.input)
			if tt.err == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			} else {
				assert.ErrorIs(t, err, tt.err)
			}
		})
	}
}
