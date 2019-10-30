package config

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
)

// matchStr check regexp match
func matchStr(rgxp string, compare string) bool {
	r, err := regexp.Compile(rgxp)
	if err != nil {
		log.Fatalf("invalid regexp: %s", rgxp)
	}
	return r.MatchString(strings.ToLower(compare))
}

// getSSHConfigPath get a default ssh user config filepath
func getSSHConfigPath(homeDir string) (string, error) {
	fp := path.Join(homeDir, "/.ssh/config")
	if ok := fileExists(fp); !ok {
		return "", errors.New("ssh user config not found: ~/.ssh/config")
	}
	return fp, nil
}

// getUserHomeDir get homedir from current user
func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
