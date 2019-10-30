package config

import (
	"bufio"
	"errors"
	"io"
	"sort"
	"strings"
)

func parseConfig(reader io.Reader) ([]Host, error) {
	var (
		hosts         []Host
		tmpHost       Host
		fileTextLines []string
		aliases       []string
		tmpMapHost    = make(map[string]Host)
	)
	tmpHost.Options = make(map[string]string)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fileTextLines = append(fileTextLines, scanner.Text())
	}

	if len(fileTextLines) == 0 {
		return nil, errors.New("error: ssh config file is empty")
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	for lineNum, scannerText := range fileTextLines {
		lineText := strings.TrimSpace(scannerText)
		if matchStr("^#", lineText) {
			// ignore comments
			continue
		}

		if matchStr("^host ", lineText) {
			if tmpHost.Alias != "" {
				aliases = append(aliases, tmpHost.Alias)
				tmpMapHost[tmpHost.Alias] = tmpHost
				tmpHost.Alias = ""
				tmpHost.Options = make(map[string]string)
			}
			// set host alias name
			alias := lineText[4:len(lineText)]
			tmpHost.Alias = strings.TrimSpace(alias)
			continue
		}

		if lineText != "" {
			// remove comments
			line := strings.Split(lineText, "#")[0]
			line = strings.TrimSpace(line)

			// split default when value is after space
			lineCfg := strings.SplitN(line, " ", 2)
			if strings.Contains(line, "=") {
				lineCfg = strings.SplitN(line, "=", 2)
			}

			if len(lineCfg) == 2 {
				tmpHost.Options[strings.ToLower(lineCfg[0])] = strings.TrimSpace(lineCfg[1])
			}
		}

		if len(fileTextLines) == lineNum+1 && tmpHost.Alias != "" {
			aliases = append(aliases, tmpHost.Alias)
			tmpMapHost[tmpHost.Alias] = tmpHost
		}
	}

	// sort results
	sort.Strings(aliases)
	for _, alias := range aliases {
		if item, ok := tmpMapHost[alias]; ok {
			hosts = append(hosts, item)
		}
	}

	return hosts, nil
}
