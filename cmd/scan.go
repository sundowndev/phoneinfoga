package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/output"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
)

type ScanCmdOptions struct {
	Number           string
	DisabledScanners []string
	PluginPaths      []string
	EnvFiles         []string
}

func init() {
	// Register command
	opts := &ScanCmdOptions{}
	cmd := NewScanCmd(opts)
	rootCmd.AddCommand(cmd)

	// Register flags
	cmd.PersistentFlags().StringVarP(&opts.Number, "number", "n", "", "The phone number to scan (E164 or international format)")
	cmd.PersistentFlags().StringArrayVarP(&opts.DisabledScanners, "disable", "D", []string{}, "Scanner to skip for this scan")
	cmd.PersistentFlags().StringArrayVar(&opts.PluginPaths, "plugin", []string{}, "Extra scanner plugin to use for the scan")
	cmd.PersistentFlags().StringSliceVar(&opts.EnvFiles, "env-file", []string{}, "Env files to parse environment variables from (looks for .env by default)")
	// scanCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "Text file containing a list of phone numbers to scan (one per line)")
	// scanCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output to save scan results")
}

func NewScanCmd(opts *ScanCmdOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "Scan a phone number",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := godotenv.Load(opts.EnvFiles...)
			if err != nil {
				logrus.WithField("error", err).Debug("Error loading .env file")
			}

			runScan(opts)
		},
	}
}

func runScan(opts *ScanCmdOptions) {
	fmt.Fprintf(color.Output, color.WhiteString("Running scan for phone number %s...\n\n"), opts.Number)

	if valid := number.IsValid(opts.Number); !valid {
		logrus.WithFields(map[string]interface{}{
			"input": opts.Number,
			"valid": valid,
		}).Debug("Input phone number is invalid")
		exitWithError(errors.New("given phone number is not valid"))
	}

	num, err := number.NewNumber(opts.Number)
	if err != nil {
		exitWithError(err)
	}

	for _, p := range opts.PluginPaths {
		err := remote.OpenPlugin(p)
		if err != nil {
			exitWithError(err)
		}
	}

	f := filter.NewEngine()
	f.AddRule(opts.DisabledScanners...)

	remoteLibrary := remote.NewLibrary(f)
	remote.InitScanners(remoteLibrary)

	// Scanner options are currently not used in CLI
	result, errs := remoteLibrary.Scan(num, remote.ScannerOptions{})

	err = output.GetOutput(output.Console, color.Output).Write(result, errs)
	if err != nil {
		exitWithError(err)
	}
}
