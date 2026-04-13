package librespot

import (
	"testing"

	"github.com/dubeyKartikay/lazyspotify/core/utils"
)

func TestAudioBackendForOS(t *testing.T) {
	tests := []struct {
		goos string
		want string
	}{
		{goos: "darwin", want: "audio-toolbox"},
		{goos: "linux", want: "alsa"},
		{goos: "freebsd", want: "alsa"},
	}

	for _, tt := range tests {
		if got := audioBackendForOS(tt.goos); got != tt.want {
			t.Fatalf("audioBackendForOS(%q) = %q, want %q", tt.goos, got, tt.want)
		}
	}
}

func TestMakeLibrespotConfigUsesConfiguredDaemonLogLevel(t *testing.T) {
	cfg := utils.AppConfig{}
	cfg.Librespot.Host = "127.0.0.1"
	cfg.Librespot.Port = 4040
	cfg.Librespot.Daemon.LogLevel = "WARN"

	got := makeLibrespotConfig(cfg, "user-id", "token")

	if got.LogLevel != "warn" {
		t.Fatalf("makeLibrespotConfig(...).LogLevel = %q, want %q", got.LogLevel, "warn")
	}
}

func TestMakeLibrespotConfigDefaultsDaemonLogLevelToError(t *testing.T) {
	cfg := utils.AppConfig{}
	cfg.Librespot.Host = "127.0.0.1"
	cfg.Librespot.Port = 4040

	got := makeLibrespotConfig(cfg, "user-id", "token")

	if got.LogLevel != "error" {
		t.Fatalf("makeLibrespotConfig(...).LogLevel = %q, want %q", got.LogLevel, "error")
	}
}
