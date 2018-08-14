package main

import (
	"github.com/urfave/cli"
)

func cliOptions() []cli.Command {

	// Commands
	cliCommands := []cli.Command{
		{
			Name:      "add",
			Usage:     "Add a new SSH alias to ~/.ssh/config",
			UsageText: "nssh add [-h] [--key SSH_KEY_FILE] name user@host:port -o option=value -o option2=value2",
			HelpName:  "add",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key",
					Usage: "SSH Key",
				},
				cli.StringSliceFlag{
					Name:  "o",
					Usage: "Option config",
				},
			},
			Action: Add,
		},
		{
			Name:      "del",
			Usage:     "Delete SSH by alias name",
			UsageText: "nssh del [-h] name",
			HelpName:  "del",
			Action:    Delete,
		},
		{
			Name:      "list",
			Usage:     "List SSH alias",
			UsageText: "nssh list [-h]",
			HelpName:  "list",
			Action:    List,
		},
		{
			Name:      "search",
			Usage:     "Search SSH alias by given search text",
			UsageText: "nssh search [-h] text",
			HelpName:  "search",
			Action:    Search,
		},
		{
			Name:      "edit",
			Usage:     "Edit SSH alias by name",
			UsageText: "nssh edit [-h] [--key SSH_KEY_FILE] name user@host:port -o option=value -o option2=value2 -p",
			HelpName:  "edit",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key",
					Usage: "SSH Key",
				},
				cli.StringFlag{
					Name:  "r",
					Usage: "Rename host",
				},
				cli.StringSliceFlag{
					Name:  "o",
					Usage: "Option config",
				},
				cli.BoolFlag{
					Name:  "p",
					Usage: "Preserve options",
				},
				cli.BoolFlag{
					Name:  "f",
					Usage: "Force edit, don't ask to preserve another options",
				},
			},
			Action: Edit,
		},
		{
			Name:      "export-csv",
			Usage:     "Export SSH alias to csv file (Only name, user, hostname and port)",
			UsageText: "nssh export [-h] -file /tmp/sshconfig.csv",
			HelpName:  "export-csv",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Usage: "Output CSV file",
				},
			},
			Action: ExportCSV,
		},
		{
			Name:      "backup",
			Usage:     "Backup SSH alias",
			UsageText: "nssh backup [-h] -file=/tmp/backup.sshconfig",
			HelpName:  "backup",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Usage: "Output backup file",
				},
			},
			Action: Backup,
		},
	}

	return cliCommands
}
