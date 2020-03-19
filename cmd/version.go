package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sundowndev/phoneinfoga/pkg/config"
	"github.com/sundowndev/phoneinfoga/pkg/utils"
)

func init() {
	// Register command
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version of the tool",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoggerService.Infoln("PhoneInfoga", config.Version)
		utils.LoggerService.Infoln("Coded by Sundowndev https://github.com/sundowndev/PhoneInfoga")
	},
}
