package actions

import (
	"bufio"
	"errors"
	"os"
	"strings"

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
	hostOptions, err := getHostOptions(alias, uri, sshkey, options)
	if err != nil {
		return err
	}

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
		if confirmProceedUpdate() {
			return errors.New("Operation cancelled")
		}
	}

	hostOptions, err := getHostOptions(oldAlias, uri, sshkey, options)
	if err != nil {
		return err
	}

	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	if err := sshConfig.UpdateHost(oldAlias, newAlias, hostOptions, preserveOptions); err != nil {
		return err
	}

	return nil
}

func confirmProceedUpdate() bool {
	reader := bufio.NewReader(os.Stdin)
	utils.Printc(utils.GlobalTitleColor, "Proceed without preserve another options? y/n\n")
	text, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	if strings.ToLower(strings.TrimSpace(text)) == "y" {
		return true
	}

	return false
}
