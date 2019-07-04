package cmd

import (
	"fmt"
	"os"

	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "touchstone",
	Short: "Touchstone is a benchmarking suite for CRI-compatible container runtimes.",
}

var version = "dev"
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("touchstone %s\n", version)

		client, err := runtime.NewClient(util.GetCRIEndpoint(cri))
		if err != nil {
			logrus.WithError(err).Fatal("could not connect")
		}
		fmt.Println(client.Version())
	},
}

var (
	cri, handler string
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(benchmarkCmd)
	rootCmd.PersistentFlags().StringVarP(&cri, "cri", "c", "containerd", "CRI runtime")
	rootCmd.PersistentFlags().StringVarP(&handler, "runtime-handler", "r", "runc", "OCI handler")
}

// Execute runs the command executor.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("failed to execute:", err)
		os.Exit(1)
	}
}
