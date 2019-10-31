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
	NewHost(string, map[string]string) error
	DeleteHost(string) error
	write() error
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
func validateNewHost(alias string, options map[string]string, bufferHosts []Host) error {
	if len(alias) == 0 {
		return errors.New("New host alias can't be empty")
	}

	if _, ok := options["hostname"]; !ok && alias != "*" {
		return errors.New("Hostname not found")
	}

	for key, value := range options {
		if len(value) == 0 {
			return fmt.Errorf("Parameter %s can't be empty", key)
		}
	}
	for _, bh := range bufferHosts {
		if bh.Alias == alias {
			return errors.New("Host already exist")
		}
	}

	return nil
}

func (cfg *sshConfig) NewHost(alias string, options map[string]string) error {
	if err := validateNewHost(alias, options, cfg.hosts); err != nil {
		return err
	}

	newHost := Host{
		Alias:   alias,
		Options: options,
	}

	// TODO: Sort after insert a new host
	cfg.hosts = append(cfg.hosts, newHost)

	return cfg.write()
}

// DeleteHost delete host configuration
// need call WriteConfig for apply
func (cfg *sshConfig) DeleteHost(alias string) error {
	for idx, x := range cfg.hosts {
		if strings.ToLower(x.Alias) == strings.ToLower(alias) {
			cfg.hosts = append(cfg.hosts[:idx], cfg.hosts[idx+1:]...)
			return cfg.write()
		}
	}

	return fmt.Errorf("Host [%s] not found", alias)
}

// WriteConfig write buffer hosts to a temporary file and
// then overwrite the ssh configuration file
func (cfg *sshConfig) write() error {
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

	if err := copyFile(swpFile.Name(), cfg.configFile); err != nil {
		return err
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
