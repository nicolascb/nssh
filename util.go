package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

func HostConnect(connection string) map[string]string {
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

func CurrentUser() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return usr.Username
}

func Contains(sub string, options []string, str ...string) bool {
	sub = strings.ToLower(sub)

	for _, s := range str {
		s = strings.ToLower(s)
		if strings.Contains(s, sub) {
			return true
		}
	}

	for _, o := range options {
		o = strings.ToLower(o)
		if strings.Contains(o, sub) {
			return true
		}
	}

	return false
}

func CopyFile(dst string) error {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Format sshconfig filepath
	sshconfig := fmt.Sprintf("%s/.ssh/config", usr.HomeDir)

	original, err := os.Open(sshconfig)

	if err != nil {
		return err
	}

	defer original.Close()

	// Create new file
	backup, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer backup.Close()

	_, err = io.Copy(backup, original)

	if err != nil {
		return err
	}

	err = backup.Sync()
	if err != nil {
		return err
	}

	return nil
}

func CreateCSV(dst string) (int, error) {

	list := CreateList()
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

	err = writer.Write(header)

	if err != nil {
		return len(list), err
	}

	for _, h := range list {
		line := []string{h.Name, h.User, h.Hostname, h.Port}
		err := writer.Write(line)
		if err != nil {
			return len(list), err
		}
	}

	return len(list), nil
}
