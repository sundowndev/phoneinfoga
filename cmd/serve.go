package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sundowndev/phoneinfoga/v2/build"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"github.com/sundowndev/phoneinfoga/v2/web"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api/handlers"
	"log"
	"net/http"
	"os"
)

type ServeCmdOptions struct {
	HttpPort         int
	DisableClient    bool
	DisabledScanners []string
	PluginPaths      []string
	EnvFiles         []string
}

func init() {
	// Register command
	opts := &ServeCmdOptions{}
	cmd := NewServeCmd(opts)
	rootCmd.AddCommand(cmd)

	// Register flags
	cmd.PersistentFlags().IntVarP(&opts.HttpPort, "port", "p", 5000, "HTTP port")
	cmd.PersistentFlags().BoolVar(&opts.DisableClient, "no-client", false, "Disable web client (REST API only)")
	cmd.PersistentFlags().StringArrayVarP(&opts.DisabledScanners, "disable", "D", []string{}, "Scanner to skip for the scans")
	cmd.PersistentFlags().StringArrayVar(&opts.PluginPaths, "plugin", []string{}, "Extra scanner plugin to use for the scans")
	cmd.PersistentFlags().StringSliceVar(&opts.EnvFiles, "env-file", []string{}, "Env files to parse environment variables from (looks for .env by default)")
}

func NewServeCmd(opts *ServeCmdOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve web client",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := godotenv.Load(opts.EnvFiles...)
			if err != nil {
				logrus.WithField("error", err).Debug("Error loading .env file")
			}

			for _, p := range opts.PluginPaths {
				err := remote.OpenPlugin(p)
				if err != nil {
					exitWithError(err)
				}
			}

			// Initialize remote library
			f := filter.NewEngine()
			f.AddRule(opts.DisabledScanners...)
			handlers.Init(f)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if build.IsRelease() && os.Getenv("GIN_MODE") == "" {
				gin.SetMode(gin.ReleaseMode)
			}

			srv, err := web.NewServer(opts.DisableClient)
			if err != nil {
				log.Fatal(err)
			}

			addr := fmt.Sprintf(":%d", opts.HttpPort)
			fmt.Printf("Listening on %s\n", addr)
			if err := srv.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		},
	}
}
