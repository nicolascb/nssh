package main

import (
	"log"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli"
)

var (
	app = cli.NewApp()
)

func main() {
	app.Name = APP_NAME
	app.Description = APP_DESCRIPTION
	app.Version = APP_VERSION
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Nicolas Barbosa",
			Email: "ndevbarbosa@gmail.com",
		},
	}

	app.Usage = APP_USAGE
	app.UsageText = APP_USAGE_TEXT

	app.Commands = CliOptions()
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
