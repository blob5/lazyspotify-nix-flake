package buildinfo

import "testing"

func TestText(t *testing.T) {
	prevVersion := Version
	prevCommit := Commit
	prevBuildDate := BuildDate
	prevPackagedDaemonPath := PackagedDaemonPath
	t.Cleanup(func() {
		Version = prevVersion
		Commit = prevCommit
		BuildDate = prevBuildDate
		PackagedDaemonPath = prevPackagedDaemonPath
	})

	Version = "1.2.3"
	Commit = "abc123"
	BuildDate = "2026-04-13T00:00:00Z"
	PackagedDaemonPath = "/usr/lib/lazyspotify/lazyspotify-librespot"

	got := Text()
	want := "version=1.2.3\ncommit=abc123\nbuild_date=2026-04-13T00:00:00Z\npackaged_daemon_path=/usr/lib/lazyspotify/lazyspotify-librespot\n"
	if got != want {
		t.Fatalf("Text() = %q, want %q", got, want)
	}
}
