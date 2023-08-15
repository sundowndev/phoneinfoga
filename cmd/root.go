package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "phoneinfoga [COMMANDS] [OPTIONS]",
	Short:   "Advanced information gathering & OSINT tool for phone numbers",
	Long:    "PhoneInfoga is one of the most advanced tools to scan phone numbers using only free resources.",
	Example: "phoneinfoga scan -n <number>",
}

// Execute is a function that executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Fprintf(color.Error, "%s\n", color.RedString(err.Error()))
	os.Exit(1)
}
