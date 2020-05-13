package cmd

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	api "gopkg.in/sundowndev/phoneinfoga.v2/api"
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

		api.Serve(router, disableClient)

		httpPort := ":" + strconv.Itoa(httpPort)

		srv := &http.Server{
			Addr:    httpPort,
			Handler: router,
		}

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	},
}
