package utils

import (
	"strings"

	"github.com/fatih/color"
)

type logger struct {
	newColor func(value ...color.Attribute) *color.Color
}

// Infoln logs an info message
func (l *logger) Infoln(s ...string) {
	l.newColor(color.FgCyan).Println("[i]", strings.Join(s, " "))
}

// Warnln logs an warning message
func (l *logger) Warnln(s ...string) {
	l.newColor(color.FgYellow).Println("[*]", strings.Join(s, " "))
}

// Errorln logs an error message
func (l *logger) Errorln(s ...string) {
	l.newColor(color.FgRed).Println("[!]", strings.Join(s, " "))
}

// Successln logs a success message
func (l *logger) Successln(s ...string) {
	l.newColor(color.FgGreen).Println("[+]", strings.Join(s, " "))
}

// Successf logs a success message
func (l *logger) Successf(format string, messages ...interface{}) {
	l.newColor(color.FgGreen).Printf("[+] "+format, messages...)
	l.newColor().Printf("\n")
}

// LoggerService is the default logger instance
var LoggerService = &logger{
	newColor: color.New,
}
