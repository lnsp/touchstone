package cmd

import (
	"fmt"
	"log"

	"github.com/lnsp/touchstone/pkg/framework"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the benchmark suite.",
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, _ := cmd.Flags().GetString("endpoint")
		_, err := framework.NewClient(endpoint)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Running benchmarks ... done!")
	},
}
