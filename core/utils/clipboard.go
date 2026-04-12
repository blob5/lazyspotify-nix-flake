package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

type clipboardCommand struct {
	name string
	args []string
}

var clipboardCommands = []clipboardCommand{
	{name: "pbcopy"},
	{name: "wl-copy"},
	{name: "xclip", args: []string{"-selection", "clipboard"}},
	{name: "xsel", args: []string{"--clipboard", "--input"}},
}

func CopyToClipboard(text string) error {
	for _, clipboardCmd := range clipboardCommands {
		path, err := exec.LookPath(clipboardCmd.name)
		if err != nil {
			continue
		}

		cmd := exec.Command(path, clipboardCmd.args...)
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no supported clipboard command found")
}
