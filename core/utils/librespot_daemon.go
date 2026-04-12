package utils

import (
	"fmt"
	"os"
)

var defaultLibrespotDaemonPath = ""

func ResolveLibrespotDaemonCmd(configured []string) ([]string, error) {
	if len(configured) > 0 {
		return configured, nil
	}

	if defaultLibrespotDaemonPath == "" {
		return nil, fmt.Errorf(
			"librespot daemon path is not configured: no packaged default was compiled in; set librespot.daemon.cmd in config or build with -ldflags \"-X github.com/dubeyKartikay/lazyspotify/core/utils.defaultLibrespotDaemonPath=...\"",
		)
	}

	if err := validateDaemonExecutable(defaultLibrespotDaemonPath); err != nil {
		return nil, fmt.Errorf(
			"librespot daemon not available at packaged default %q: %w; install lazyspotify-librespot there or set librespot.daemon.cmd in config",
			defaultLibrespotDaemonPath,
			err,
		)
	}

	return []string{defaultLibrespotDaemonPath}, nil
}

func validateDaemonExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory")
	}
	if info.Mode()&0o111 == 0 {
		return fmt.Errorf("path is not executable")
	}
	return nil
}
