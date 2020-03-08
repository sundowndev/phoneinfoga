package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is the current version of the tool
const Version = "v1.0.0"

var number string
var input string  // TODO
var output string // TODO

var rootCmd = &cobra.Command{
	Use:   "phoneinfoga [COMMANDS] [OPTIONS]",
	Short: "Advanced information gathering & OSINT tool for phone numbers",
	Long:  `PhoneInfoga is one of the most advanced tools to scan phone numbers using only free resources.`,
}

// Execute is a function that executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
