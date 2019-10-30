// Package config provides interface for manipulating ssh configuration file
package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type SSHConfig interface {
	DeleteHost(string) error
	WriteConfig() error
}

type sshConfig struct {
	configFile string
	hosts      []Host
}

func LoadUserConfig() (SSHConfig, error) {
	// Get current user homedir
	homeDir, err := getUserHomeDir()
	if err != nil {
		return nil, err
	}

	// Get if exist ~/.ssh/config file
	configFile, err := getSSHConfigPath(homeDir)
	if err != nil {
		return nil, err
	}

	// ~/.ssh/config exist, proceed to open and parse
	// hosts
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	hosts, err := parseConfig(f)
	if err != nil {
		return nil, err
	}

	return &sshConfig{
		configFile: configFile,
		hosts:      hosts,
	}, nil
}

// DeleteHost delete host configuration
// need call WriteConfig for apply
func (cfg *sshConfig) DeleteHost(alias string) error {
	for idx, x := range cfg.hosts {
		if strings.ToLower(x.Alias) == strings.ToLower(alias) {
			cfg.hosts = append(cfg.hosts[:idx], cfg.hosts[idx+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Host [%s] not found", alias)
}

// WriteConfig write buffer hosts to a temporary file and
// then overwrite the ssh configuration file
func (cfg *sshConfig) WriteConfig() error {
	swpFile, err := ioutil.TempFile("/tmp", "config.*.swp")
	if err != nil {
		return err
	}

	defer os.Remove(swpFile.Name())

	for _, hostData := range cfg.hosts {
		textConfig, err := hostDecoder(hostData)
		if err != nil {
			return err
		}
		if _, err := swpFile.WriteString(textConfig); err != nil {
			return err
		}
	}

	return nil
}

func hostDecoder(host Host) (string, error) {
	if host.Alias == "" {
		return "", errors.New("Invalid alias")
	}

	var (
		header  string
		options string
	)
	header = fmt.Sprintf("Host %s\n", host.Alias)
	for key, value := range host.Options {
		// TODO: Arrumar padding
		options += fmt.Sprintf("%5s %s\n", key, value)
	}

	return header + options, nil

}
