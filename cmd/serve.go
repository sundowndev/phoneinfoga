package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	api "github.com/sundowndev/phoneinfoga/api"
)

var httpPort int
var disableClient bool

func init() {
	// Register command
	rootCmd.AddCommand(serveCmd)

	// Register flags
	serveCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 5000, "HTTP port")
	serveCmd.PersistentFlags().BoolVar(&disableClient, "no-client", false, "Disable web client (REST API only)")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve web client",
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()

		api.Serve(router, httpPort, disableClient)
	},
}
