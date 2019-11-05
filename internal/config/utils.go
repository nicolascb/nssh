package config

import (
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"unicode"
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
		newConfigFile, err := os.Create(fp)
		if err != nil {
			return "", err
		}

		defer newConfigFile.Close()
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

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func tabToSpace(input string) string {
	var result []string

	for _, i := range input {
		switch {
		// all these considered as space, including tab \t
		// '\t', '\n', '\v', '\f', '\r',' ', 0x85, 0xA0
		case unicode.IsSpace(i):
			result = append(result, " ") // replace tab with space
		case !unicode.IsSpace(i):
			result = append(result, string(i))
		}
	}
	return strings.Join(result, "")
}
