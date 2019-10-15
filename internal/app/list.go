package app

import (
	"fmt"
	"strings"

	"github.com/nicolascb/nssh/internal/utils"
	"github.com/nicolascb/nsshconfig"
)

// SSHHost represent a ~/.ssh/config entry
type SSHHost struct {
	Name           string
	User           string
	Hostname       string
	Port           string
	AnotherOptions []string
}

// GetSSHEntries return a ~/.ssh/config parsed
func GetSSHEntries() ([]SSHHost, error) {
	var (
		entries []SSHHost
		host    *SSHHost
	)

	// Parse sshconfig
	if err := nsshconfig.LoadConfig(); err != nil {
		return nil, err
	}

	general, _ := nsshconfig.GetEntryByHost("*")

	for _, e := range nsshconfig.Hosts() {
		host = new(SSHHost)
		if e.Host != "*" {
			host.Name = e.Host

			if general != nil {
				for gk, gv := range general.Options {
					switch strings.ToLower(gk) {
					case "hostname":
						host.Hostname = gv
					case "user":
						host.User = gv
					case "port":
						host.Port = gv
					}
				}
			}

			for ek, ev := range e.Options {
				switch strings.ToLower(ek) {
				case "hostname":
					host.Hostname = ev
				case "user":
					host.User = ev
				case "port":
					host.Port = ev
				default:
					host.AnotherOptions = append(host.AnotherOptions, fmt.Sprintf("%s: %s", ek, ev))
				}
			}

			if err := host.setDefaults(); err != nil {
				return nil, err
			}

			entries = append(entries, *host)
		}

	}

	return entries, nil
}

func printList(list []SSHHost) {
	By(Prop("Name", true)).Sort(list)
	utils.Printc(utils.DefaultColor, "List:\n")
	for _, x := range list {
		str := fmt.Sprintf(" -> %s@%s:%s\n", x.User, x.Hostname, x.Port)

		utils.Printc(utils.TitleColor, fmt.Sprintf("	%s", x.Name))
		fmt.Print(str)

		for a, c := range x.AnotherOptions {
			if a == 0 {
				fmt.Printf("	  [options]")
				fmt.Printf("  %s\n", c)
				continue
			}

			fmt.Printf("		    %s\n", c)
		}
	}
}

func searchAlias(substr string) ([]SSHHost, error) {
	entries, err := GetSSHEntries()
	if err != nil {
		return nil, err
	}

	var foundEntries []SSHHost
	for _, x := range entries {
		if utils.Contains(substr, x.AnotherOptions, x.Hostname, x.Name, x.Port, x.User) {
			foundEntries = append(foundEntries, x)
		}
	}

	return foundEntries, err
}

func (host *SSHHost) setDefaults() error {
	notData := "not_specified"

	if host.Name == "" {
		host.Name = notData
	}

	if host.User == "" {
		username, err := utils.GetCurrentUser()
		if err != nil {
			return err
		}
		host.User = username
	}

	if host.Hostname == "" {
		host.Hostname = notData
	}

	if host.Port == "" {
		host.Port = "22"
	}

	return nil
}
