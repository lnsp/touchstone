package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/benchmark/suites"
	"github.com/lnsp/touchstone/pkg/config"
	"github.com/lnsp/touchstone/pkg/visual"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	pattern    string
	outDir     string
	visualFile string
)

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run the benchmark suite",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := filepath.Glob(pattern)
		if err != nil {
			logrus.WithError(err).Fatal("failed expand glob")
		}
		var (
			index   = benchmark.NewIndex()
			entries []benchmark.MatrixEntry
		)
		for _, file := range files {
			logrus.WithField("file", file).Info("loading benchmark file")
			cfg, err := config.Parse(file)
			if err != nil {
				logrus.WithError(err).Fatal("failed parse config")
			}
			out, err := cfg.MapOutput(outDir)
			if err != nil {
				logrus.WithError(err).Fatal("failed map output")
			}
			defer out.Close()
			matrix, err := cfg.Matrix()
			if err != nil {
				logrus.WithError(err).Fatal("failed build matrix")
			}
			results, err := matrix.Run()
			if err != nil {
				logrus.WithError(err).Fatal("failed matrix run")
			}
			encoder := json.NewEncoder(out)
			if err := encoder.Encode(results); err != nil {
				logrus.WithError(err).Fatal("failed json encode")
			}
			entries = append(entries, results...)
			// update index
			matrix.Index(index)
		}
		if err := visual.Write(outDir+visualFile, entries, index); err != nil {
			logrus.WithError(err).Fatal("failed write")
		}
	},
}

var listFilter []string
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available benchmarks",
	Run: func(cmd *cobra.Command, args []string) {
		filtered := suites.All()
		for _, filter := range listFilter {
			filtered = benchmark.Filter(suites.All(), filter)
		}
		for _, b := range filtered {
			fmt.Println(b.Name())
		}
	},
}

func init() {
	benchmarkCmd.Flags().StringVarP(&pattern, "file", "f", "default.yaml", "Input benchmark configuration")
	benchmarkCmd.Flags().StringVarP(&outDir, "dir", "d", "", "Output destination directory")
	benchmarkCmd.Flags().StringVarP(&visualFile, "html-file", "x", "index.html", "HTML visualisation file name")
	listCmd.Flags().StringSliceVarP(&listFilter, "filter", "f", nil, "Filter expression")
}
