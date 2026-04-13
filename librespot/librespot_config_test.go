package librespot

import "testing"

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
