package benchmark

import (
	"fmt"
	"strings"

	"github.com/lnsp/touchstone/pkg/runtime"
)

type Benchmark interface {
	Name() string
	Run(client *runtime.Client, handler string) (interface{}, error)
}

type BenchmarkSuite []Benchmark

func (bs BenchmarkSuite) Name() string {
	names := make([]string, len(bs))
	for i := range bs {
		names[i] = bs[i].Name()
	}
	return fmt.Sprintf("BenchmarkSuite [%s]", strings.Join(names, ", "))
}

func (bs BenchmarkSuite) Run(client *runtime.Client, handler string) (interface{}, error) {
	reports := make([]struct {
		Name   string      `json:"name"`
		Report interface{} `json:"report"`
	}, len(bs))
	for i := range bs {
		report, err := bs[i].Run(client, handler)
		if err != nil {
			return nil, fmt.Errorf("failed to run suite %s: %v", bs[i].Name(), err)
		}
		reports[i].Name = bs[i].Name()
		reports[i].Report = report
	}
	return reports, nil
}

// getBenchmarkID computes a unique string for this benchmark.
func getBenchmarkID(b Benchmark) string {
	return fmt.Sprintf("Benchmark-%s-%s", b.Name(), runtime.NewUUID())
}
