package utils

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

// GetCurrentUser return current username
func GetCurrentUser() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.Username, nil
}

// Contains check slice string contain a string
func Contains(sub string, options []string, str ...string) bool {
	sub = strings.ToLower(sub)

	for _, s := range str {
		if strings.Contains(strings.ToLower(s), sub) {
			return true
		}
	}

	for _, o := range options {
		if strings.Contains(strings.ToLower(o), sub) {
			return true
		}
	}

	return false
}

// CopySSHConfigFile copy ~/.ssh/config
func CopySSHConfigFile(dst string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	// Format sshconfig filepath
	sshconfig := fmt.Sprintf("%s/.ssh/config", usr.HomeDir)

	original, err := os.Open(sshconfig)

	if err != nil {
		return err
	}

	defer original.Close()

	// Create new file
	backup, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer backup.Close()

	if _, err = io.Copy(backup, original); err != nil {
		return err
	}

	if err = backup.Sync(); err != nil {
		return err
	}

	return nil
}
