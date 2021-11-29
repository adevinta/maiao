package system_test

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/adevinta/maiao/pkg/system"
)

func TestGetEnvOrDefault(t *testing.T) {
	t.Cleanup(system.Reset)
	os.Setenv("some-key", "some-value")

	assert.Equal(t, system.GetenvOrDefault("other-key", "default-value"), "default-value")
	assert.Equal(t, system.GetenvOrDefault("some-key", "default-value"), "some-value")
}

func TestEnvConsidersCurrentEnvironmentVariables(t *testing.T) {
	t.Cleanup(system.Reset)
	_, ok := system.Env()["some-key"]
	assert.False(t, ok)
	os.Setenv("some-key", "some-value")
	assert.Equal(t, "some-value", system.Env()["some-key"])
}

func TestResetRecoversDeletedEnvironmentVariables(t *testing.T) {
	env := system.Env()

	os.Unsetenv("PATH")
	require.NotEqual(t, env, system.Env())
	system.Reset()
	assert.Equal(t, env, system.Env())
}

func TestResetRestoresModifiedEnvironmentVariables(t *testing.T) {
	env := system.Env()

	os.Setenv("PATH", "alternative-value")
	require.NotEqual(t, env, system.Env())
	system.Reset()
	assert.Equal(t, env, system.Env())
}

func TestResetRemovesAddedEnvironmentVariables(t *testing.T) {
	env := system.Env()

	os.Setenv("some-key", "some-value")
	require.NotEqual(t, env, system.Env())
	system.Reset()
	assert.Equal(t, env, system.Env())
}

func TestResetRestoresDefaultFileSystem(t *testing.T) {
	fs := system.DefaultFileSystem
	system.DefaultFileSystem = afero.NewMemMapFs()
	require.NotEqual(t, fs, system.DefaultFileSystem)
	system.Reset()
	assert.Equal(t, fs, system.DefaultFileSystem)
}

func TestResetRestoresWorkingDirectory(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir("../"))
	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.NotEqual(t, wd, cwd)
	system.Reset()
	cwd, err = os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, wd, cwd)
}
