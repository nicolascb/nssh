package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nicolascb/nssh/internal/actions"
	"github.com/nicolascb/nssh/internal/config"
	"github.com/nicolascb/nssh/internal/utils"
	"github.com/urfave/cli"
)

// Delete alias
func Delete(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowCommandHelpAndExit(c, "del", 1)
	}

	alias := c.Args().First()
	if err := actions.Delete(alias); err != nil {
		return err
	}

	utils.OkColor.Printf("Successfully deleted alias [%s]\n", alias)
	return nil
}

// Add alias
func Add(c *cli.Context) error {
	// Args validation
	if c.NArg() < 2 && c.Args().First() != config.GeneralDefinitions {
		cli.ShowCommandHelpAndExit(c, "add", 1)
	}

	var (
		alias   string
		sshkey  string
		uri     string
		options []string
	)

	// Host name
	alias = c.Args().First()
	// SSH Key
	sshkey = c.String("key")
	// Custom options
	options = c.StringSlice("o")

	if len(c.Args()) > 1 {
		uri = c.Args()[1]
	}

	if err := actions.Add(alias, uri, sshkey, options); err != nil {
		return err
	}

	// OK
	utils.OkColor.Printf("Successfully added alias [%s]\n", c.Args().First())
	return nil
}

// Edit alias
func Edit(c *cli.Context) error {
	// Valid Arguments
	if c.NArg() < 2 && !c.IsSet("o") && !c.IsSet("r") && !c.IsSet("key") {
		cli.ShowCommandHelpAndExit(c, "edit", 1)
	}

	var (
		oldAlias string
		newAlias string
		sshkey   string
		uri      string
		options  []string
	)

	// Old host name
	oldAlias = c.Args().First()
	// New host name
	newAlias = c.String("r")
	// SSH Key
	sshkey = c.String("key")
	// Custom options
	options = c.StringSlice("o")

	if oldAlias == config.GeneralDefinitions && newAlias != "" {
		return errors.New("(*) General can not be renamed")
	}

	// Check connection is passed
	if c.NArg() == 2 && oldAlias != config.GeneralDefinitions {
		uri = c.Args()[1]
	}

	if err := actions.Edit(oldAlias, newAlias, uri, sshkey, options, c.Bool("f"), c.Bool("p")); err != nil {
		return err
	}

	utils.OkColor.Printf("Successfully edited [%s]\n", c.Args().First())
	return nil
}

// List aliases in ~/.ssh/config
func List(ct *cli.Context) error {
	textOutput, err := actions.Print()
	if err != nil {
		return err
	}

	fmt.Print(textOutput)
	return nil
}

// Search alias by text
func Search(c *cli.Context) error {
	// Check args
	if c.NArg() < 1 {
		cli.ShowCommandHelpAndExit(c, "search", 1)
	}

	word := c.Args().First()
	textOutput, err := actions.Search(word)
	if err != nil {
		return err
	}

	fmt.Print(textOutput)
	return nil
}

// Backup sshconfig
func Backup(c *cli.Context) error {
	// File to backup
	file := strings.TrimSpace(c.String("file"))
	if file == "" {
		cli.ShowCommandHelpAndExit(c, "backup", 1)
	}

	// Copy backup
	if err := utils.CopySSHConfigFile(file); err != nil {
		return err
	}

	// OK
	// utils.Printc(utils.OkColor, fmt.Sprintf("Finished backup [%s]", file))
	return nil
}

// ExportCSV save sshconfig in csv file
func ExportCSV(c *cli.Context) error {
	// Output file
	file := strings.TrimSpace(c.String("file"))
	if file == "" {
		cli.ShowCommandHelpAndExit(c, "export-csv", 1)
	}

	// List and create CSV
	rows, err := generateCSV(file)
	if err != nil {
		return err
	}
	print(rows)
	// CSV OK
	// utils.Printc(utils.OkColor, fmt.Sprintf("Finished export csv [%s] %d aliases", file, rows))
	return nil
}
