// Package config provides interface for manipulating ssh configuration file
package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// GerneralDefinitions identifies the host that contains parameters for all connections
	GeneralDefinitions = "*"
)

const (
	// ActionTypeNew id to identify new hosts
	ActionTypeNew = iota + 1
	// ActionTypeUpdate id to identify updates
	ActionTypeUpdate
)

// SSHConfig interface to manipulate ssh settings
type SSHConfig interface {
	NewHost(string, map[string]string) error
	UpdateHost(string, string, map[string]string, bool) error
	DeleteHost(string) error
	write() error
}

type sshConfig struct {
	configFile string
	hosts      []Host
}

// LoadUserConfig read sshconfig user file
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

func validateOptions(actionType int, needHostname bool, alias string, options map[string]string, bufferHosts []Host) error {
	if len(alias) == 0 {
		return errors.New("Host alias can't be empty")
	}

	if needHostname && alias != GeneralDefinitions {
		if _, ok := options["hostname"]; !ok {
			return errors.New("Hostname not found")
		}
	}

	for key, value := range options {
		if len(value) == 0 {
			return fmt.Errorf("Parameter %s can't be empty", key)
		}
	}

	for _, bh := range bufferHosts {
		if actionType == ActionTypeNew && bh.Alias == alias {
			return errors.New("Host already exist")
		}

		if bh.Alias == alias {
			return nil
		}
	}

	if actionType == ActionTypeUpdate {
		return fmt.Errorf("Update failed: host [%s] not found", alias)
	}

	return nil
}

func (cfg *sshConfig) NewHost(alias string, options map[string]string) error {
	if err := validateOptions(ActionTypeNew, true, alias, options, cfg.hosts); err != nil {
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
		if strings.EqualFold(x.Alias, alias) {
			cfg.hosts = append(cfg.hosts[:idx], cfg.hosts[idx+1:]...)
			return cfg.write()
		}
	}

	return fmt.Errorf("Host [%s] not found", alias)
}

func (cfg *sshConfig) UpdateHost(oldAlias, newAlias string, options map[string]string, preserveOptions bool) error {
	var needHostname bool
	if !preserveOptions {
		needHostname = true
	}

	if err := validateOptions(ActionTypeUpdate, needHostname, oldAlias, options, cfg.hosts); err != nil {
		return err
	}

	hostAlias := oldAlias
	if len(newAlias) > 0 {
		hostAlias = newAlias
	}

	editedHost := Host{
		Alias:   hostAlias,
		Options: options,
	}

	for index, host := range cfg.hosts {
		if host.Alias == oldAlias {
			if !preserveOptions {
				cfg.hosts[index] = editedHost
				break
			}

			for key, val := range editedHost.Options {
				cfg.hosts[index].Options[key] = val
			}

			cfg.hosts[index].Alias = hostAlias
			break
		}
	}

	return cfg.write()
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
