package actions

import (
	"errors"
	"fmt"
	"os"

	"github.com/nicolascb/nssh/internal/config"
	"github.com/nicolascb/nssh/internal/utils"
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

func Print() (int, error) {
	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return 0, err
	}

	var general config.Host
	hosts := sshConfig.Hosts()
	for _, host := range hosts {
		if host.Alias == config.GeneralDefinitions {
			general = host
			continue
		}

		uri := host.Options["hostname"]
		if port, ok := host.Options["port"]; ok {
			uri = fmt.Sprintf("%s:%s", uri, port)
		}

		if user, ok := host.Options["user"]; ok {
			uri = fmt.Sprintf("%s@%s", user, uri)
		}

		output := fmt.Sprintf(" -> %s\n", uri)
		utils.TitleColor.Printf("	%s", host.Alias)
		fmt.Print(output)
	}

	utils.DefaultColor.Printf("general: %v\n", general)
	return len(hosts), nil
}
