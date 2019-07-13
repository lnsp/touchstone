package cmd

import (
	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/benchmark/suites"
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
			CRIs:  []string{"containerd", "crio"},
			OCIs:  []string{"runc", "runsc"},
			Items: suites.All,
			Runs:  10,
		}
		if err := matrix.Run(outfile); err != nil {
			logrus.WithError(err).Error("error while running benchmark")
		}
	},
}

func init() {
	benchmarkCmd.Flags().StringVarP(&stdout, "output", "o", "-", "Output target for reports")
}
