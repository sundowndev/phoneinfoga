package cmd

import (
	"github.com/spf13/cobra"
	api "github.com/sundowndev/phoneinfoga/api"
)

var httpPort int

func init() {
	// Register command
	rootCmd.AddCommand(serveCmd)

	// Register flags
	serveCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 5000, "HTTP port")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve web client",
	Run: func(cmd *cobra.Command, args []string) {
		api.Serve(httpPort)
	},
}
