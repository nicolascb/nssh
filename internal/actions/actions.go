package actions

import (
	"errors"
	"os"
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

// Add host to sshconfig file
func Add(alias, uri, sshkey string, options []string) error {
	hostOptions := getHostOptions(alias, uri, sshkey, options)
	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	if err := sshConfig.NewHost(alias, hostOptions); err != nil {
		return err
	}

	return nil
}

// Edit host
func Edit(oldAlias, newAlias, uri, sshkey string, options []string, forceUpdate, preserveOptions bool) error {
	if !preserveOptions && !forceUpdate {
		if !confirmProceedUpdate(os.Stdin) {
			return errors.New("Operation cancelled")
		}
	}

	hostOptions := getHostOptions(oldAlias, uri, sshkey, options)
	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	if err := sshConfig.UpdateHost(oldAlias, newAlias, hostOptions, preserveOptions); err != nil {
		return err
	}

	return nil
}

// Print return hosts in output text format
func Print() (string, error) {
	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return "", err
	}

	output := formatOutput(sshConfig.Hosts())
	return output, nil
}

// Search host by text
// Compare if host/hostname match searched word
func Search(word string) (string, error) {
	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return "", err
	}

	var foundHosts []config.Host
	hosts := sshConfig.Hosts()
	for _, host := range hosts {
		if strings.Contains(host.Alias, word) {
			foundHosts = append(foundHosts, host)
			continue
		}

		if hostname, ok := host.Options["hostname"]; ok {
			if strings.Contains(hostname, word) {
				foundHosts = append(foundHosts, host)
				continue
			}
		}
	}

	output := formatOutput(foundHosts)
	return output, nil
}
