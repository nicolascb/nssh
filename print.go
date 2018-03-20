package main

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	defaultMessage     = color.New(color.FgWhite, color.Bold)
	globalTitleMessage = color.New(color.FgYellow, color.Bold)

	titleMessage = color.New(color.FgBlue, color.Bold)
	// globalKeyStyle = color.New(color.FgCyan, color.Bold)

	okMessage  = color.New(color.FgGreen)
	errMessage = color.New(color.FgRed)
)

func PrintList(le []ListEntry) {
	By(Prop("Name", true)).Sort(le)
	defaultMessage.Printf("List:\n")
	for _, x := range le {
		str := fmt.Sprintf(" -> %s@%s:%s\n", x.User, x.Hostname, x.Port)
		titleMessage.Printf("	%s", x.Name)
		fmt.Printf(str)

		for a, c := range x.AnotherOptions {
			if a == 0 {
				fmt.Printf("	  [options]")
				fmt.Printf("  %s\n", c)
			} else {
				fmt.Printf("		    %s\n", c)
			}
		}
	}
	return
}

func PrintGeneral(general map[string]string) {
	globalTitleMessage.Println("	(*) General Options")

	for i, g := range general {
		fmt.Printf("		%s: %s\n", i, g)
	}
}

func PrintOK(msg string) {
	okMessage.Printf("✔ OK: ")
	defaultMessage.Printf("%s\n", msg)
}

func PrintErr(err error) {
	if err != nil {
		errMessage.Printf("✗ ERROR: ")
		defaultMessage.Printf("%s\n", err.Error())
	}
}
