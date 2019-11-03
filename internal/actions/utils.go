package actions

import (
	"bufio"
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

func getHostOptions(alias, uri, sshkey string, options []string) (map[string]string, error) {
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
			return nil, errors.New("Hostname not found")
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

	return hostOptions, nil
}
