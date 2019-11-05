package actions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/nicolascb/nssh/internal/config"
	"github.com/nicolascb/nssh/internal/utils"
)

func parseURI(uri string) (string, string, string) {
	var (
		user     string
		port     string
		hostname string
	)

	tmpURI := strings.Split(uri, "@")
	if len(tmpURI) > 1 {
		user = tmpURI[0]
		hostname = tmpURI[1]
	}

	if len(hostname) == 0 {
		hostname = uri
	}

	tmpPort := strings.Split(hostname, ":")
	if len(tmpPort) > 1 {
		hostname = tmpPort[0]
		port = tmpPort[1]
	}

	return user, port, hostname
}

func getHostOptions(alias, uri, sshkey string, options []string) map[string]string {
	hostOptions := make(map[string]string)

	if alias != config.GeneralDefinitions {
		user, port, hostname := parseURI(uri)

		if len(user) > 0 {
			hostOptions["user"] = user
		}

		if len(port) > 0 {
			hostOptions["port"] = port
		}

		if len(hostname) > 0 {
			hostOptions["hostname"] = hostname
		}
	}

	if len(sshkey) > 0 {
		hostOptions["identityfile"] = sshkey
	}

	for _, opt := range options {
		splited := strings.Split(opt, "=")
		if len(splited) > 1 {
			key := strings.ToLower(splited[0])
			val := strings.ToLower(splited[1])
			hostOptions[key] = val
		}
	}

	return hostOptions
}

func confirmProceedUpdate(out io.Reader) bool {
	reader := bufio.NewReader(out)
	utils.GlobalTitleColor.Print("Proceed without preserve another options? y/n\n")

	text, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	if strings.ToLower(strings.TrimSpace(text)) == "y" {
		return true
	}

	return false
}

func formatOutput(hosts []config.Host) string {
	var (
		buf     bytes.Buffer
		general config.Host
	)

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

		// Print host alias
		utils.TitleColor.Fprintf(&buf, "	%s", host.Alias)

		// Print host connection
		output := fmt.Sprintf(" -> %s\n", uri)
		fmt.Fprint(&buf, output)

		// Print host options
		for key, val := range host.Options {
			if isPrintableOption(key) {
				fmt.Fprintf(&buf, "		%s: %s\n", key, val)
			}
		}
	}

	if general.Options != nil {
		utils.GlobalTitleColor.Fprintf(&buf, "	(*) General options")
		for key, val := range general.Options {
			if isPrintableOption(key) {
				fmt.Fprintf(&buf, "		%s: %s\n", key, val)
			}
		}
	}

	utils.DefaultColor.Fprintf(&buf, "Found %d entries\n", len(hosts))
	return buf.String()
}

func isPrintableOption(option string) bool {
	reservedOptions := []string{"user", "port", "hostname"}
	for _, rOption := range reservedOptions {
		if strings.EqualFold(rOption, option) {
			return false
		}
	}

	return true
}
