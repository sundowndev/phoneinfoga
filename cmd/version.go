package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/sundowndev/phoneinfoga.v2/pkg/config"
	"gopkg.in/sundowndev/phoneinfoga.v2/pkg/utils"
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
