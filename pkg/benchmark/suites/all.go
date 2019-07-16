package suites

import "github.com/lnsp/touchstone/pkg/benchmark"

func All() (All []benchmark.Benchmark) {
	All = append(All, Performance...)
	All = append(All, Operations...)
	All = append(All, Scalability...)
	return All
}
