package buildinfo

import (
	"fmt"
	"io"
)

var (
	Version            = "dev"
	Commit             = "unknown"
	BuildDate          = "unknown"
	PackagedDaemonPath = ""
)

func Text() string {
	return fmt.Sprintf(
		"version=%s\ncommit=%s\nbuild_date=%s\npackaged_daemon_path=%s\n",
		Version,
		Commit,
		BuildDate,
		PackagedDaemonPath,
	)
}

func PrintVersion(w io.Writer) error {
	_, err := io.WriteString(w, Text())
	return err
}
