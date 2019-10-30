package actions

import "github.com/nicolascb/nssh/internal/config"

// Delete remove host alias
func Delete(alias string) error {

	sshConfig, err := config.LoadUserConfig()
	if err != nil {
		return err
	}

	if err := sshConfig.DeleteHost(alias); err != nil {
		return err
	}

	if err := sshConfig.WriteConfig(); err != nil {
		return err
	}

	return nil
}
