package main

import (
	"os"

	"github.com/nicolascb/nsshconfig"
)

func ParseConfig() {
	err := nsshconfig.LoadConfig()
	if err != nil {
		PrintErr(err)
		os.Exit(1)
	}
}
