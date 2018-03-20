package main

import (
	"fmt"
	"strings"

	"github.com/nicolascb/nsshconfig"
)

type ListEntry struct {
	Name           string
	User           string
	Hostname       string
	Port           string
	AnotherOptions []string
}

var (
	hL []ListEntry
	h  ListEntry
)

func CreateList() []ListEntry {

	ParseConfig()

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

			h = SetDefaults(h)
			hL = append(hL, h)
		}

	}

	return hL
}

func SetDefaults(e ListEntry) ListEntry {
	notData := "not_specified"

	if e.Name == "" {
		e.Name = notData
	}

	if e.User == "" {
		e.User = CurrentUser()
	}

	if e.Hostname == "" {
		e.Hostname = notData
	}

	if e.Port == "" {
		e.Port = "22"
	}

	return e
}
