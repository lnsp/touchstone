package cmd

import (
	"fmt"
	"os"

	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Version = "dev"

var knownCRIs = []string{"containerd", "crio"}
var rootCmd = &cobra.Command{
	Use:   "touchstone",
	Short: "Touchstone is a benchmarking suite for CRI-compatible container runtimes.",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("touchstone %s\n", Version)
		for _, cri := range knownCRIs {
			client, err := runtime.NewClient(util.GetCRIEndpoint(cri))
			if err != nil {
				logrus.WithError(err).WithField("cri", cri).Error("failed connect")
			}
			fmt.Println(client.Name(), client.Version())
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(benchmarkCmd)
	rootCmd.AddCommand(listCmd)
}

// Execute runs the command executor.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
