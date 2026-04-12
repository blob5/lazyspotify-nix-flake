package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveLibrespotDaemonCmd_ConfigOverrideWins(t *testing.T) {
	prev := defaultLibrespotDaemonPath
	defaultLibrespotDaemonPath = ""
	t.Cleanup(func() {
		defaultLibrespotDaemonPath = prev
	})

	configured := []string{"/custom/librespot", "--flag"}
	got, err := ResolveLibrespotDaemonCmd(configured)
	if err != nil {
		t.Fatalf("ResolveLibrespotDaemonCmd returned error: %v", err)
	}
	if len(got) != len(configured) {
		t.Fatalf("got %d args, want %d", len(got), len(configured))
	}
	for i := range configured {
		if got[i] != configured[i] {
			t.Fatalf("got arg %d = %q, want %q", i, got[i], configured[i])
		}
	}
}

func TestResolveLibrespotDaemonCmd_UsesCompiledDefault(t *testing.T) {
	path := makeExecutable(t)

	prev := defaultLibrespotDaemonPath
	defaultLibrespotDaemonPath = path
	t.Cleanup(func() {
		defaultLibrespotDaemonPath = prev
	})

	got, err := ResolveLibrespotDaemonCmd(nil)
	if err != nil {
		t.Fatalf("ResolveLibrespotDaemonCmd returned error: %v", err)
	}
	if len(got) != 1 || got[0] != path {
		t.Fatalf("got %v, want [%q]", got, path)
	}
}

func TestResolveLibrespotDaemonCmd_EmptyCompiledDefaultFails(t *testing.T) {
	prev := defaultLibrespotDaemonPath
	defaultLibrespotDaemonPath = ""
	t.Cleanup(func() {
		defaultLibrespotDaemonPath = prev
	})

	_, err := ResolveLibrespotDaemonCmd(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no packaged default was compiled in") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "librespot.daemon.cmd") {
		t.Fatalf("expected config guidance in error: %v", err)
	}
}

func TestResolveLibrespotDaemonCmd_MissingCompiledDefaultFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-librespot")

	prev := defaultLibrespotDaemonPath
	defaultLibrespotDaemonPath = path
	t.Cleanup(func() {
		defaultLibrespotDaemonPath = prev
	})

	_, err := ResolveLibrespotDaemonCmd(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), path) {
		t.Fatalf("expected path in error, got: %v", err)
	}
}

func TestResolveLibrespotDaemonCmd_NonExecutableCompiledDefaultFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "librespot")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	prev := defaultLibrespotDaemonPath
	defaultLibrespotDaemonPath = path
	t.Cleanup(func() {
		defaultLibrespotDaemonPath = prev
	})

	_, err := ResolveLibrespotDaemonCmd(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not executable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func makeExecutable(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "librespot")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	return path
}
