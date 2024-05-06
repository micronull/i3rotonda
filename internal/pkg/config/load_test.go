package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/micronull/i3rotonda/internal/pkg/config"
)

func TestLoad(t *testing.T) {
	tmpHomePath, err := os.MkdirTemp("", "i3rotonda_test")
	assert.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tmpHomePath) })

	t.Setenv("HOME", tmpHomePath)

	tmpCfgPath, err := os.UserConfigDir()
	assert.NoError(t, err)

	tmpCfgPath = filepath.Join(tmpCfgPath, "i3rotonda")
	tmpCfg := filepath.Join(tmpCfgPath, "config.yml")

	err = os.MkdirAll(tmpCfgPath, os.ModePerm)
	assert.NoError(t, err)

	data := `debug: true
workspaces:
  exclude:
    - some
    - some2
    - some3
`

	err = os.WriteFile(tmpCfg, []byte(data), 0644)
	assert.NoError(t, err)

	got, err := config.Load()
	assert.NoError(t, err)

	expected := config.Config{
		Debug: true,
	}
	expected.Workspaces.Exclude = []string{"some", "some2", "some3"}

	assert.Equal(t, expected, got)
}
