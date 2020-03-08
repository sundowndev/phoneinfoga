package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type logger struct{}

// Infoln logs an info message
func (l *logger) Infoln(s ...string) {
	color.New(color.FgCyan).Println("[i]", strings.Join(s, " "))
}

// Warnln logs an warning message
func (l *logger) Warnln(s ...string) {
	color.New(color.FgYellow).Println("[*]", strings.Join(s, " "))
}

// Errorln logs an error message
func (l *logger) Errorln(s ...string) {
	color.New(color.FgRed).Println("[!]", strings.Join(s, " "))
}

// Successln logs a success message
func (l *logger) Successln(s ...string) {
	color.New(color.FgGreen).Println("[+]", strings.Join(s, " "))
}

// Successf logs a success message
func (l *logger) Successf(format string, a ...interface{}) {
	color.New(color.FgGreen).Printf("[+] "+format, a...)
	fmt.Printf("\n")
}

// LoggerService is the default logger instance
var LoggerService = &logger{}
