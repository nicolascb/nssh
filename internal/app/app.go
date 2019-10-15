package app

import (
	"encoding/csv"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli"
)

// NewNsshApp create cli interface
func NewNsshApp() *cli.App {
	app := cli.NewApp()
	app.Name = appName
	app.Description = appDescription
	app.Version = appVersion
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Nicolas Barbosa",
			Email: "ndevbarbosa@gmail.com",
		},
	}

	app.Usage = appUsage
	app.UsageText = appUsageText

	app.Commands = cliOptions()
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	return app
}

func parseHostConnection(connection string) map[string]string {
	options := make(map[string]string)
	if strings.Contains(connection, "@") {
		options["user"] = strings.Split(connection, "@")[0]
		options["hostname"] = strings.Split(connection, "@")[1]
	}

	if strings.Contains(connection, ":") {
		options["port"] = strings.Split(connection, ":")[1]
		if _, ok := options["hostname"]; ok {
			options["hostname"] = strings.Split(options["hostname"], ":")[0]
		} else {
			options["hostname"] = strings.Split(connection, ":")[0]
		}
	}

	if _, ok := options["hostname"]; !ok {
		options["hostname"] = connection
	}

	return options
}

func generateCSV(dst string) (int, error) {

	list, err := GetSSHEntries()
	if err != nil {
		return 0, err
	}

	file, err := os.Create(dst)
	if err != nil {
		return len(list), err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Alias",
		"User",
		"Hostname",
		"Port"}

	if err = writer.Write(header); err != nil {
		return len(list), err
	}

	for _, h := range list {
		line := []string{h.Name, h.User, h.Hostname, h.Port}
		if err := writer.Write(line); err != nil {
			return len(list), err
		}
	}

	return len(list), nil
}
