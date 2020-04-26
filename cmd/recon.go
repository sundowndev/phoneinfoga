package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/sundowndev/phoneinfoga.v2/pkg/utils"
)

func init() {
	// Register command
	rootCmd.AddCommand(reconCmd)

	// Register flags
	reconCmd.PersistentFlags().StringVarP(&number, "number", "n", "", "The phone number to scan (E164 or international format)")
	reconCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "Text file containing a list of phone numbers to scan (one per line)")
	reconCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output to save scan results")
}

var reconCmd = &cobra.Command{
	Deprecated: "Not implemented yet.",
	Use:        "recon",
	Short:      "Launch custom format reconnaissance",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoggerService.Infoln("Custom recon for phone number", number)

		utils.LoggerService.Errorln("Not implemented yet.")

		utils.LoggerService.Infoln("Job finished.")
	},
}
