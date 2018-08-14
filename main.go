package main

import (
	"os"
	"sort"
	"time"

	"github.com/urfave/cli"
)

var (
	app = cli.NewApp()
)

func main() {
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
	if err := app.Run(os.Args); err != nil {
		printErr(err)
	}
}
