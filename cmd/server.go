package cmd

import (
	"net/http"

	"github.com/Faione/perf_exporter/server"
	"github.com/spf13/cobra"
)

var (
	keyListenAddr  = "listen.addr"
	keyMetricPath  = "metric.path"
	keyMaxRequests = "max.requests"
)

func init() {

	serveCmd := &cobra.Command{
		Use:   "serve [options]",
		Short: "run on server mode, collect perf event and expose as metric",
		RunE:  serve,
	}

	flags := serveCmd.Flags()

	flags.StringP(keyListenAddr, "l", ":9900", "Address on which to expose metrics and web interface")
	flags.String(keyMetricPath, "/metrics", "Path under which to expose metrics.")
	flags.Int(
		keyMaxRequests,
		0,
		"Maximum number of parallel scrape requests. Use 0 to disable.",
	)
	rootCmd.AddCommand(serveCmd)
}

func serve(c *cobra.Command, args []string) error {

	addr, err := c.Flags().GetString(keyListenAddr)
	if err != nil {
		return err
	}

	maxRequests, err := c.Flags().GetInt(keyMaxRequests)
	if err != nil {
		return err
	}

	metricPath, err := c.Flags().GetString(keyMetricPath)
	if err != nil {
		return err
	}

	server := server.NewServer(&server.Config{
		MetricPath:        metricPath,
		MetricMaxRequests: maxRequests,
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: server,
	}

	return srv.ListenAndServe()
}
