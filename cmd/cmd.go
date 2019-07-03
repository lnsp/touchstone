package cmd

import (
	"fmt"
	"os"

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
