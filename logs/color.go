// Copyright 2017 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logs

import (
	"os"

	isatty "github.com/mattn/go-isatty"
)

const escape = "\x1b"

// Foreground text colors.
const (
	fgRed    = "31"
	fgGreen  = "32"
	fgYellow = "33"
)

var noColor = !isatty.IsTerminal(os.Stdout.Fd())

func colorize(code, text string) string {
	if noColor {
		return text
	}
	return format(code) + text + unformat()
}

func format(code string) string {
	return escape + "[" + code + "m"
}

func unformat() string {
	return escape + "[0m"
}

func RedString(text string) string {
	return colorize(fgRed, text)
}

func GreenString(text string) string {
	return colorize(fgGreen, text)
}

func YellowString(text string) string {
	return colorize(fgYellow, text)
}
