package main

import (
	"fmt"
	"os"

	"github.com/nicolascb/nssh/internal/app"
	"github.com/nicolascb/nssh/internal/utils"
)

func main() {
	cmd := app.NewNsshApp()
	if err := cmd.Run(os.Args); err != nil {
		utils.Printc(utils.ErrColor, fmt.Sprintf("failed: %v\n", err))
	}
}
