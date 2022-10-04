package aesgcm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

func promptPassPhrase(promptScript string) (string, error) {
	if promptScript == "" {
		return consolePrompt()
	}

	buf := bytes.NewBuffer([]byte{})
	cmd := exec.Command(promptScript)
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func consolePrompt() (string, error) {
	fmt.Fprint(os.Stderr, "Enter passphrase:")
	text, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return strings.TrimSpace(string(text)), nil
}
