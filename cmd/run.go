package cmd

import (
	"fmt"

	"github.com/lnsp/touchstone/pkg/framework"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the benchmark suite.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := framework.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Running benchmarks ... done!")
	},
}
