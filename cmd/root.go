package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Faione/perf_exporter/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          filepath.Base(os.Args[0]) + " [server | stat]",
		Long:         "collect perf event and expose as metric",
		SilenceUsage: true,
		Version:      version.Version.String(),
	}
)

func init() {
	globalFlags := rootCmd.Flags()

	_ = globalFlags
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(1)
}
