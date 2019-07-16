package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Output the report index of the given benchmarks",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := filepath.Glob(pattern)
		if err != nil {
			logrus.WithError(err).Fatal("failed expand glob")
		}
		index := benchmark.NewIndex()
		for _, file := range files {
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
			matrix.Index(index)
		}
		indexJSON, err := json.Marshal(index)
		if err != nil {
			logrus.WithError(err).Fatal("failed marshal index")
		}
		fmt.Println(string(indexJSON))
	},
}

func init() {
	indexCmd.Flags().StringVarP(&pattern, "file", "f", "default.yaml", "Input benchmark configuration")
}
