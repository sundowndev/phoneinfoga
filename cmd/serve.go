package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sundowndev/phoneinfoga/v2/web"
	"log"
	"net/http"
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
		srv, err := web.NewServer(disableClient)
		if err != nil {
			log.Fatal(err)
		}

		addr := fmt.Sprintf(":%d", httpPort)

		fmt.Printf("Listening on %s\n", addr)
		if err := srv.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	},
}
