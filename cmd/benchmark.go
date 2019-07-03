package cmd

import (
	"encoding/json"

	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/runtime"
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

		client, err := runtime.NewClient(util.GetCRIEndpoint(cri))
		if err != nil {
			logrus.WithError(err).Fatal("failed to create client")
		}
		defer client.Close()
		btt := []struct {
			Name     string
			Attempts int
			Suite    benchmark.Benchmark
		}{
			{"StartupTime", 10, benchmark.SuiteStartupTime},
		}
		encoder := json.NewEncoder(outfile)
		for _, btc := range btt {
			logrus.WithFields(logrus.Fields{
				"Name":     btc.Name,
				"Attempts": btc.Attempts,
			}).Info("running benchmark suite")
			for i := 0; i < btc.Attempts; i++ {
				report, err := btc.Suite.Run(client, handler)
				if err != nil {
					logrus.WithError(err).Error("error while running benchmark")
				}
				if err := encoder.Encode(report); err != nil {
					logrus.WithError(err).Fatal("failed to marshal report")
				}
			}
		}
	},
}

func init() {
	benchmarkCmd.Flags().StringVarP(&stdout, "output", "o", "-", "Output target for reports")
}
