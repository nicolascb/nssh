package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nicolascb/nssh/internal/utils"
	"github.com/nicolascb/nsshconfig"
	"github.com/urfave/cli"
)

// Delete alias
func Delete(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowCommandHelpAndExit(c, "del", 1)
	}

	host := c.Args().First()

	// Parse sshconfig
	if err := nsshconfig.LoadConfig(); err != nil {
		return err
	}

	// Check exist host
	if _, err := nsshconfig.GetEntryByHost(host); err != nil {
		return fmt.Errorf("Host [%s] not found", host)
	}

	// Host exist, proceed delete:
	err := nsshconfig.Delete(host)
	if err != nil {
		return err
	}

	// Write file
	err = nsshconfig.WriteConfig()
	if err != nil {
		return err
	}

	// Host deleted
	utils.Printc(utils.OkColor, fmt.Sprintf("Successfully deleted host [%s]", host))
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
		optsMap = parseHostConnection(connection)
	}

	// Parse sshconfig
	if err := nsshconfig.LoadConfig(); err != nil {
		return err
	}

	// Check if exist host
	if nsshconfig.ExistHost(name) {
		// Host already exist, print error and exit
		return fmt.Errorf("Host [%s] already exist", name)
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
	if err := nsshconfig.New(name, optsMap); err != nil {
		return err
	}

	// Write file
	if err := nsshconfig.WriteConfig(); err != nil {
		return err
	}

	// OK
	utils.Printc(utils.OkColor, fmt.Sprintf("Successfully added host [%s]\n", c.Args().First()))
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
		return fmt.Errorf("(*) General can not be renamed")
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
		optsMap = parseHostConnection(connection)
	}

	// Parse sshconfig
	if err := nsshconfig.LoadConfig(); err != nil {
		return err
	}

	// Get host
	host, err := nsshconfig.GetEntryByHost(name)

	// Host not found
	if err != nil {
		return fmt.Errorf("Host [%s] not found", name)
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
			utils.Printc(utils.GlobalTitleColor, "Proceed without preserve another options? y/n\n")
			text, _ := reader.ReadString('\n')
			text = strings.ToLower(strings.TrimSpace(text))
			if text != "y" {
				return fmt.Errorf("Operation cancelled")
			}
		}
		if name != "*" {
			if _, ok := optsMap["hostname"]; !ok {
				return fmt.Errorf("Hostname not especified: use -p to preserve options")
			}
		}
		host.Options = optsMap
	} else {
		for k, v := range optsMap {
			host.Options[k] = v
		}
	}

	// Save host
	if err = host.Save(); err != nil {
		return err
	}

	// Write file
	if err = nsshconfig.WriteConfig(); err != nil {
		return err
	}

	// OK
	utils.Printc(utils.OkColor, fmt.Sprintf("Successfully edited [%s]", c.Args().First()))
	return nil
}

// List aliases in ~/.ssh/config
func List(ct *cli.Context) error {
	// Create alias list
	list, err := GetSSHEntries()
	if err != nil {
		return err
	}

	// Get general options
	general, _ := nsshconfig.GetEntryByHost("*")

	if len(list) > 0 {
		printList(list)
	}

	// If general exist, print general:
	if general != nil {
		utils.Printc(utils.GlobalTitleColor, "	(*) General Options")

		for i, g := range general.Options {
			fmt.Printf("		%s: %s\n", i, g)
		}
	}

	// Default message, found alias
	utils.Printc(utils.DefaultColor, fmt.Sprintf("\nFound %d entries\n", len(list)))
	return nil
}

// Search alias by text
func Search(c *cli.Context) error {
	// Check args
	if c.NArg() < 1 {
		cli.ShowCommandHelpAndExit(c, "search", 1)
	}
	// Exec search
	found, err := searchAlias(c.Args().First())
	if err != nil {
		return err
	}

	// Alias not found
	if len(found) == 0 {
		utils.Printc(utils.DefaultColor, fmt.Sprintf("No matches found for [%s]\n", c.Args().First()))
		return nil
	}

	// Print
	printList(found)
	utils.Printc(utils.DefaultColor, fmt.Sprintf("\nFound %d entries.\n", len(found)))
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
	utils.Printc(utils.OkColor, fmt.Sprintf("Finished backup [%s]", file))
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

	// CSV OK
	utils.Printc(utils.OkColor, fmt.Sprintf("Finished export csv [%s] %d aliases", file, rows))
	return nil
}
