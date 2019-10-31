package actions

import (
	"errors"
	"strings"

	"github.com/nicolascb/nssh/internal/config"
)

// Delete remove host alias
func Delete(alias string) error {

	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	if err := sshConfig.DeleteHost(alias); err != nil {
		return err
	}

	return nil
}

func Add(alias, uri, sshkey string, options []string) error {
	hostOptions := make(map[string]string)

	if alias != "*" {
		user, port, hostname := parseURI(uri)

		if len(user) > 0 {
			hostOptions["user"] = user
		}

		if len(port) > 0 {
			hostOptions["port"] = port
		}

		if len(hostname) > 0 {
			hostOptions["hostname"] = hostname
		} else {
			return errors.New("Hostname not found")
		}
	}

	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	for _, opt := range options {
		splited := strings.Split(opt, "=")
		if len(splited) > 1 {
			key := strings.ToLower(splited[0])
			val := strings.ToLower(splited[1])
			hostOptions[key] = val
		}
	}

	if len(sshkey) > 0 {
		hostOptions["identityfile"] = sshkey
	}

	if err := sshConfig.NewHost(alias, hostOptions); err != nil {
		return err
	}

	return nil
}
