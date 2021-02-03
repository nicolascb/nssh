package app

import (
	"sort"
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
