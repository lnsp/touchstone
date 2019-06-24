package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the benchmark suite.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running benchmarks ... done!")
	},
}
