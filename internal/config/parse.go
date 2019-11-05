package config

import (
	"bufio"
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

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	for lineNum, scannerText := range fileTextLines {
		lineText := strings.TrimSpace(scannerText)
		if matchStr("^#", lineText) {
			// ignore comments
			continue
		}

		if matchStr("^host ", tabToSpace(lineText)) {
			if tmpHost.Alias != "" {
				aliases = append(aliases, tmpHost.Alias)
				tmpMapHost[tmpHost.Alias] = tmpHost
				tmpHost.Alias = ""
				tmpHost.Options = make(map[string]string)
			}
			// set host alias name
			alias := lineText[4:]
			tmpHost.Alias = strings.TrimSpace(tabToSpace(alias))
			continue
		}

		if lineText != "" {
			// remove comments
			line := strings.Split(lineText, "#")[0]
			line = strings.TrimSpace(line)

			// split default when value is after space
			lineCfg := strings.Split(tabToSpace(line), " ")
			if strings.Contains(line, "=") {
				lineCfg = strings.Split(line, "=")
			}
			if len(lineCfg) > 1 {

				key := lineCfg[0]
				val := strings.Join(lineCfg[1:], "")
				tmpHost.Options[strings.ToLower(key)] = strings.TrimSpace(val)
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
