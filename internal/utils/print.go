package utils

import (
	"github.com/fatih/color"
)

var (
	DefaultColor     = color.New(color.FgWhite, color.Bold)
	GlobalTitleColor = color.New(color.FgYellow, color.Bold)
	TitleColor       = color.New(color.FgBlue, color.Bold)
	OkColor          = color.New(color.FgGreen)
	ErrColor         = color.New(color.FgRed)
)

func Printc(color *color.Color, msg string) {
	color.Print(msg)
}
