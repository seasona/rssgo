package test

import (
	"testing"

	"github.com/seasona/rssgo/internal"
)

func TestLoadConfig(t *testing.T) {
	url := "../config.json"

	config := internal.Config{}

	config.LoadConfig(url)

	t.Run("keySwitch", func(t *testing.T) {
		if config.KeySwitchWindows != "Tab" {
			t.Errorf("got %v, want %v", config.KeySwitchWindows, "Tab")
		}
	})

	t.Run("keyQuit", func(t *testing.T) {
		if config.KeyQuit != "Esc" {
			t.Errorf("got %v, want %v", config.KeyQuit, "Esc")
		}
	})
}
