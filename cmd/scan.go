package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/output"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"os"
)

var disabledScanners []string

func init() {
	// Register command
	rootCmd.AddCommand(scanCmd)

	// Register flags
	scanCmd.PersistentFlags().StringVarP(&inputNumber, "number", "n", "", "The phone number to scan (E164 or international format)")
	scanCmd.PersistentFlags().StringArrayVarP(&disabledScanners, "disable", "D", []string{}, "A list of scanners to skip for this scan.")
	// scanCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "Text file containing a list of phone numbers to scan (one per line)")
	// scanCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output to save scan results")
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a phone number",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runScan()
	},
}

func runScan() {
	fmt.Printf(color.WhiteString("Running scan for phone number %s...\n\n"), inputNumber)

	if valid := number.IsValid(inputNumber); !valid {
		logrus.WithFields(map[string]interface{}{
			"input": inputNumber,
			"valid": valid,
		}).Debug("Input phone number is invalid")
		fmt.Println(color.RedString("Given phone number is not valid"))
		os.Exit(1)
	}

	num, err := number.NewNumber(inputNumber)
	if err != nil {
		fmt.Println(color.RedString(err.Error()))
		os.Exit(1)
	}

	f := filter.NewEngine()
	f.AddRule(disabledScanners...)

	remoteLibrary := remote.NewLibrary(f)
	remote.InitScanners(remoteLibrary)

	result, errs := remoteLibrary.Scan(num)

	err = output.GetOutput(output.Console, os.Stdout).Write(result, errs)
	if err != nil {
		fmt.Println(color.RedString(err.Error()))
		os.Exit(1)
	}
}
