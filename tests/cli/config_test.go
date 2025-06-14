package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLI(t *testing.T) {
	skipIntegration(t)

	tempDir, err := os.MkdirTemp("", "")
	require.NoErrorf(t, err, "error creating temporary directory: %s", err)

	binPath := filepath.Join(tempDir, "act3")
	buildArgs := []string{"build", "-o", binPath, "../.."}

	c := exec.Command("go", buildArgs...)
	err = c.Run()
	require.NoErrorf(t, err, "error building binary: %s", err)

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("couldn't clean up temporary directory (%s): %s", binPath, err)
		}
	}()

	//-------------//
	//  SUCCESSES  //
	//-------------//

	t.Run("Validating correct config works", func(t *testing.T) {
		// GIVEN
		c := exec.Command(
			binPath,
			"config",
			"validate",
			"-c",
			"assets/config-good.yml",
		)

		// WHEN
		err := c.Run()

		// THEN
		assert.NoError(t, err)
	})

	t.Run("Sample config is valid", func(t *testing.T) {
		// GIVEN
		c := exec.Command(
			binPath,
			"config",
			"validate",
			"-c",
			"../../internal/cmd/assets/sample-config.yml",
		)

		// WHEN
		err := c.Run()

		// THEN
		assert.NoError(t, err)
	})

	//-------------//
	//  FAILURES   //
	//-------------//

	t.Run("Validating invalid yaml fails", func(t *testing.T) {
		// GIVEN
		c := exec.Command(
			binPath,
			"config",
			"validate",
			"-c",
			"assets/config-invalid-yaml.yml",
		)

		// WHEN
		err := c.Run()

		// THEN
		assert.Error(t, err)
	})

	t.Run("Validating invalid config fails", func(t *testing.T) {
		// GIVEN
		c := exec.Command(
			binPath,
			"config",
			"validate",
			"-c",
			"assets/config-invalid.yml",
		)

		// WHEN
		output, err := c.CombinedOutput()

		// THEN
		require.Error(t, err)
		expected := `
Error: config is not valid:
- workflow at index 1 has errors: [workflow ID is empty, repo name is invalid]
- workflow at index 2 has errors: [workflow name is empty]
- workflow at index 3 has errors: [workflow key is empty]
- workflow at index 4 has errors: [URL is invalid]
`
		assert.Contains(t, string(output), strings.TrimSpace(expected))
	})
}
