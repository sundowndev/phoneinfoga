package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/sundowndev/phoneinfoga.v2/scanners"
	"gopkg.in/sundowndev/phoneinfoga.v2/utils"
)

func init() {
	// Register command
	rootCmd.AddCommand(scanCmd)

	// Register flags
	scanCmd.PersistentFlags().StringVarP(&number, "number", "n", "", "The phone number to scan (E164 or international format)")
	// scanCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "Text file containing a list of phone numbers to scan (one per line)")
	// scanCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output to save scan results")
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a phone number",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoggerService.Infoln("Scanning phone number", number)

		if valid := utils.IsValid(number); !valid {
			utils.LoggerService.Errorln("Number is not valid.")
			os.Exit(1)
		}

		scanners.ScanCLI(number)

		utils.LoggerService.Infoln("Job finished.")
	},
}
