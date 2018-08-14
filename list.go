package main

import (
	"fmt"
	"strings"

	"github.com/nicolascb/nsshconfig"
)

// ListEntry list module
type ListEntry struct {
	Name           string
	User           string
	Hostname       string
	Port           string
	AnotherOptions []string
}

func getList() ([]ListEntry, error) {
	hL := []ListEntry{}
	h := ListEntry{}

	// Parse sshconfig
	if err := nsshconfig.LoadConfig(); err != nil {
		return hL, err
	}

	general, _ := nsshconfig.GetEntryByHost("*")

	for _, e := range nsshconfig.Hosts() {
		h = ListEntry{}
		if e.Host != "*" {
			h.Name = e.Host

			if general != nil {
				for gk, gv := range general.Options {
					switch strings.ToLower(gk) {
					case "hostname":
						h.Hostname = gv
					case "user":
						h.User = gv
					case "port":
						h.Port = gv
					}
				}
			}

			for ek, ev := range e.Options {
				switch strings.ToLower(ek) {
				case "hostname":
					h.Hostname = ev
				case "user":
					h.User = ev
				case "port":
					h.Port = ev
				default:
					h.AnotherOptions = append(h.AnotherOptions, fmt.Sprintf("%s: %s", ek, ev))
				}
			}

			h, err := setDefaults(h)
			if err != nil {
				return hL, err
			}
			hL = append(hL, h)
		}

	}

	return hL, nil
}

func setDefaults(e ListEntry) (ListEntry, error) {
	notData := "not_specified"

	if e.Name == "" {
		e.Name = notData
	}

	if e.User == "" {
		username, err := getUser()
		if err != nil {
			return e, err
		}
		e.User = username
	}

	if e.Hostname == "" {
		e.Hostname = notData
	}

	if e.Port == "" {
		e.Port = "22"
	}

	return e, nil
}
