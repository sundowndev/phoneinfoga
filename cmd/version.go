package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/sundowndev/phoneinfoga.v2/config"
	"gopkg.in/sundowndev/phoneinfoga.v2/utils"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version of the tool",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoggerService.Infoln("PhoneInfoga", config.Version, fmt.Sprintf("(%s)", config.Commit))
		utils.LoggerService.Infoln("Maintained by sundowndev https://github.com/sundowndev/PhoneInfoga")
	},
}
