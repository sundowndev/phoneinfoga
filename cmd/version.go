package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sundowndev/phoneinfoga/v2/build"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version of the tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("PhoneInfoga %s\n", build.String())
	},
}
