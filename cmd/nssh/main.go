package main

import (
	"os"

	"github.com/nicolascb/nssh/internal/app"
	"github.com/nicolascb/nssh/internal/utils"
)

func main() {
	cmd := app.NewNsshApp()
	if err := cmd.Run(os.Args); err != nil {
		utils.ErrColor.Printf("ERROR: %v\n", err)
	}
}
