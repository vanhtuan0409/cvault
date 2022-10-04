package aesgcm

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

func promptPassPhrase(prompt string) (string, error) {
	switch {
	case prompt == "":
		return consolePrompt()
	case strings.HasPrefix(prompt, "script://"):
		return scriptPrompt(strings.TrimPrefix(prompt, "script://"))
	case strings.HasPrefix(prompt, "file://"):
		return filePrompt(strings.TrimPrefix(prompt, "file://"))
	default:
		return "", errors.New("unknown prompt")
	}
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

func scriptPrompt(script string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	cmd := exec.Command(script)
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func filePrompt(p string) (string, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
