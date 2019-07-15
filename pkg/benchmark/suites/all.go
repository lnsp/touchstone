package suites

import "github.com/lnsp/touchstone/pkg/benchmark"

var All []benchmark.Benchmark

func init() {
	All = append(All, Performance...)
	All = append(All, Operations...)
}
