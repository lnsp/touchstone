package cmd

import (
	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stdout string

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run the benchmark suite.",
	Run: func(cmd *cobra.Command, args []string) {
		outfile := util.GetOutputTarget(stdout)
		defer outfile.Close()
		matrix := &benchmark.Matrix{
			CRIs:            []string{"containerd", "crio"},
			RuntimeHandlers: []string{"runc", "runsc"},
			Suite: &benchmark.Suite{
				&benchmark.Attempt{10, benchmark.SuiteStartupTime},
				&benchmark.Attempt{1, benchmark.SuitePerformance},
			},
		}
		if err := matrix.Run(outfile); err != nil {
			logrus.WithError(err).Error("error while running benchmark")
		}
	},
}

func init() {
	benchmarkCmd.Flags().StringVarP(&stdout, "output", "o", "-", "Output target for reports")
}
