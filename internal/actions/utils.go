package actions

import "strings"

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

	tmpPort := strings.Split(hostname, ":")
	if len(tmpPort) > 1 {
		hostname = tmpPort[0]
		port = tmpPort[1]
	}

	if len(hostname) == 0 {
		hostname = uri
	}

	return user, port, hostname
}
