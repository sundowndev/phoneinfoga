package utils

import (
	"strings"

	"github.com/fatih/color"
)

// Logger allows you to log messages in the terminal
type Logger struct {
	NewColor func(value ...color.Attribute) colorLogger
}

type colorLogger interface {
	Println(a ...interface{}) (int, error)
	Printf(format string, a ...interface{}) (int, error)
}

// Infoln logs an info message
func (l *Logger) Infoln(s ...string) {
	l.NewColor(color.FgCyan).Println("[i]", strings.Join(s, " "))
}

// Warnln logs an warning message
func (l *Logger) Warnln(s ...string) {
	l.NewColor(color.FgYellow).Println("[*]", strings.Join(s, " "))
}

// Errorln logs an error message
func (l *Logger) Errorln(s ...string) {
	l.NewColor(color.FgRed).Println("[!]", strings.Join(s, " "))
}

// Successln logs a success message
func (l *Logger) Successln(s ...string) {
	l.NewColor(color.FgGreen).Println("[+]", strings.Join(s, " "))
}

// Successf logs a success message
func (l *Logger) Successf(format string, messages ...interface{}) {
	l.NewColor(color.FgGreen).Printf("[+] "+format, messages...)
	l.NewColor().Printf("\n")
}

// LoggerService is the default logger instance
var LoggerService = &Logger{
	NewColor: func(value ...color.Attribute) colorLogger {
		return color.New(value...)
	},
}
