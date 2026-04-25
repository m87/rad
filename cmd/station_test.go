package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadStations(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set("stations", map[string]string{
		"lofi": "https://example.com/lofi",
	})

	stations := loadStations()
	assert.NotNil(t, stations)
	assert.Contains(t, stations, "lofi")
	assert.Equal(t, Station{
		Alias: "lofi",
		URL:   "https://example.com/lofi",
	}, stations["lofi"])
}

func TestSaveStation(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	configPath := filepath.Join(t.TempDir(), "rad.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte("{}\n"), 0o644))

	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	require.NoError(t, saveStation(Station{Alias: "jazz", URL: "https://example.com/jazz"}))
	require.NoError(t, saveStation(Station{Alias: "rock", URL: "https://example.com/rock"}))

	viper.Reset()
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	stations := loadStations()
	assert.Equal(t, "https://example.com/jazz", stations["jazz"].URL)
	assert.Equal(t, "jazz", stations["jazz"].Alias)
	assert.Equal(t, "https://example.com/rock", stations["rock"].URL)
	assert.Equal(t, "rock", stations["rock"].Alias)
}
