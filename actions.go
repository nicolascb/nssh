package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nicolascb/nsshconfig"
	"github.com/urfave/cli"
)

// Delete alias
func Delete(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowCommandHelpAndExit(c, "del", 1)
	}

	host := c.Args().First()

	// Parse ssh alias
	ParseConfig()

	// Check exist host
	if _, err := nsshconfig.GetEntryByHost(host); err != nil {
		PrintErr(fmt.Errorf("Host [%s] not found", host))
		return nil
	}

	// Host exist, proceed delete:
	err := nsshconfig.Delete(host)
	if err != nil {
		PrintErr(err)
		return nil
	}

	// Write file
	err = nsshconfig.WriteConfig()
	if err != nil {
		PrintErr(err)
		return nil
	}

	// Host deleted
	PrintOK(fmt.Sprintf("Successfully deleted host [%s]", host))
	return nil
}

// Add alias
func Add(c *cli.Context) error {
	// Args validation
	if c.NArg() < 2 && c.Args().First() != "*" {
		cli.ShowCommandHelpAndExit(c, "add", 1)
	}

	// Host name
	name := c.Args().First()
	// SSH Key
	sshkey := c.String("key")
	// Custom options
	opts := c.StringSlice("o")
	optsMap := make(map[string]string)

	if name != "*" {
		// Connection URI
		connection := c.Args()[1]
		optsMap = HostConnect(connection)
	}

	// Load config
	ParseConfig()

	// Check if exist host
	_, err := nsshconfig.GetEntryByHost(name)

	// Host already exist, print error and exit
	if err == nil {
		PrintErr(fmt.Errorf("Host [%s] already exist", name))
		return nil
	}

	// Loop option (-o)
	if len(opts) > 0 {
		for _, o := range opts {
			if strings.Contains(o, "=") {
				key := strings.Split(o, "=")[0]
				val := strings.Split(o, "=")[1]
				optsMap[strings.ToLower(key)] = val
			}
		}
	}

	// sshkey flag
	if sshkey != "" {
		optsMap["identityfile"] = sshkey
	}

	// New
	err = nsshconfig.New(name, optsMap)

	if err != nil {
		PrintErr(err)
		return nil
	}

	// Write file
	err = nsshconfig.WriteConfig()
	if err != nil {
		PrintErr(err)
		return nil
	}

	// OK
	PrintOK(fmt.Sprintf("Successfully added host [%s]", c.Args().First()))
	return nil
}

// Edit alias
func Edit(c *cli.Context) error {
	// Valid Arguments
	if c.NArg() < 2 && !c.IsSet("o") && !c.IsSet("r") && !c.IsSet("key") {
		cli.ShowCommandHelpAndExit(c, "edit", 1)
	}

	// Host name
	name := c.Args().First()
	// New name
	newname := c.String("r")

	// Check rename (*) general
	if name == "*" && newname != "" {
		PrintErr(fmt.Errorf("(*) General can not be renamed"))
		return nil
	}

	// SSH Key
	sshkey := c.String("key")
	// Custom options
	opts := c.StringSlice("o")
	// Map custom options
	optsMap := make(map[string]string)

	// Check connection is passed
	if c.NArg() == 2 && name != "*" {
		connection := c.Args()[1]
		optsMap = HostConnect(connection)
	}

	// Parse sshconfig
	ParseConfig()
	// Get host
	host, err := nsshconfig.GetEntryByHost(name)

	// Host not found
	if err != nil {
		PrintErr(fmt.Errorf("Host [%s] not found", name))
		return nil
	}

	// Rename
	if newname != "" {
		host.Host = newname
	}

	// Loop options
	if len(opts) > 0 {
		for _, o := range opts {
			if strings.Contains(o, "=") {
				key := strings.Split(o, "=")[0]
				val := strings.Split(o, "=")[1]
				optsMap[strings.ToLower(key)] = val
			}
		}
	}

	// sshkey flag
	if sshkey != "" {
		optsMap["identityfile"] = sshkey
	}

	// Preserve options
	if !c.Bool("p") {
		if !c.Bool("f") {
			reader := bufio.NewReader(os.Stdin)
			globalTitleMessage.Printf("Proceed without preserve another options? y/n\n")
			text, _ := reader.ReadString('\n')
			text = strings.ToLower(strings.TrimSpace(text))
			if text != "y" {
				PrintErr(fmt.Errorf("Operation cancelled"))
				return nil
			}
		}
		if name != "*" {
			if _, ok := optsMap["hostname"]; !ok {
				errHostname := errors.New("Hostname not especified: use -p to preserve options")
				PrintErr(errHostname)
				return nil
			}
		}
		host.Options = optsMap
	} else {
		for k, v := range optsMap {
			host.Options[k] = v
		}
	}

	// Save host
	err = host.Save()

	if err != nil {
		PrintErr(err)
		return nil
	}

	// Write file
	err = nsshconfig.WriteConfig()
	if err != nil {
		PrintErr(err)
		return nil
	}

	// OK
	PrintOK(fmt.Sprintf("Successfully edited [%s]", c.Args().First()))
	return nil
}

// List aliases in ~/.ssh/config
func List(ct *cli.Context) error {
	// Create alias list
	list := CreateList()
	// Get general options
	general, _ := nsshconfig.GetEntryByHost("*")

	if len(list) > 0 {
		// Print List
		PrintList(list)
	}

	// If general exist, print general:
	if general != nil {
		PrintGeneral(general.Options)
	}

	// Default message, found alias
	defaultMessage.Printf("\nFound %d entries\n", len(list))
	return nil
}

// Search alias by text
func Search(c *cli.Context) error {
	// Check args
	if c.NArg() < 1 {
		cli.ShowCommandHelpAndExit(c, "search", 1)
	}
	// Exec search
	found := ExecSearch(c.Args().First())

	// Alias not found
	if len(found) == 0 {
		defaultMessage.Printf("No matches found for [%s]\n", c.Args().First())
		return nil
	}

	// Print
	PrintList(foundEntries)
	defaultMessage.Printf("\nFound %d entries.\n", len(foundEntries))
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
	err := CopyFile(file)
	if err != nil {
		PrintErr(err)
		return nil
	}

	// OK
	PrintOK(fmt.Sprintf("Finished backup [%s]", file))
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
	rows, err := CreateCSV(file)
	if err != nil {
		PrintErr(err)
		return nil
	}

	// CSV OK
	PrintOK(fmt.Sprintf("Finished export csv [%s] %d aliases", file, rows))
	return nil
}
