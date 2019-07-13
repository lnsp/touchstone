package benchmark

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
	"github.com/sirupsen/logrus"
)

type Benchmark interface {
	Name() string
	Run(client *runtime.Client, handler string) (Report, error)
}

type Report interface {
	Aggregate(r Report) Report
	Scale(n int) Report
}

type ValueReport map[string]float64

func (report ValueReport) Aggregate(other Report) Report {
	if report == nil {
		return other
	}
	otherReport := other.(ValueReport)
	result := make(map[string]float64)
	for k := range report {
		result[k] = report[k] + otherReport[k]
	}
	return ValueReport(result)
}

func (report ValueReport) Scale(n int) Report {
	result := make(map[string]float64)
	for k := range report {
		result[k] = report[k] / float64(n)
	}
	return ValueReport(result)
}

// Filter operates in-place on a slice of benchmarks.
func Filter(items []Benchmark, f string) []Benchmark {
	i := 0
	for j := range items {
		if strings.HasPrefix(items[j].Name(), f) {
			items[i] = items[j]
			i++
		}
	}
	return items[:i]
}

type sumError []error

func (e sumError) Error() string {
	s := ""
	for _, err := range e {
		s += err.Error() + "\n"
	}
	return s
}

type Matrix struct {
	CRIs  []string
	OCIs  []string
	Items []Benchmark
	Runs  int
}

type MatrixEntry struct {
	CRI     string         `json:"CRI"`
	OCI     string         `json:"OCI"`
	Results []MatrixResult `json:"Results"`
}

type MatrixResult struct {
	Name       string   `json:"Name"`
	Aggregated Report   `json:"Aggregated"`
	Reports    []Report `json:"Reports"`
}

func (m *Matrix) createEntry(cri string, handler string) (*MatrixEntry, error) {
	logrus.WithFields(logrus.Fields{
		"cri":     cri,
		"handler": handler,
	}).Info("evaluating matrix entry")
	client, err := runtime.NewClient(util.GetCRIEndpoint(cri))
	if err != nil {
		return nil, fmt.Errorf("[%s:%s] failed to initialize client: %v", cri, handler, err)
	}
	defer client.Close()
	errs := sumError(nil)
	results := make([]MatrixResult, 0, len(m.Items))
	for _, bm := range m.Items {
		logrus.WithFields(logrus.Fields{
			"name": bm.Name(),
		}).Info("running benchmark")
		aggregated := Report(nil)
		reports := make([]Report, 0, m.Runs)
		for i := 0; i < m.Runs; i++ {
			report, err := m.Items[i].Run(client, handler)
			if err != nil {
				errs = append(errs, fmt.Errorf("[%s:%s] failed to run benchmark: %v", cri, handler, err))
				break
			}
			reports = append(reports, report)
			aggregated = aggregated.Aggregate(report)
		}
		results = append(results, MatrixResult{
			Name:       bm.Name(),
			Aggregated: aggregated,
			Reports:    reports,
		})
	}
	return &MatrixEntry{
		CRI:     cri,
		OCI:     handler,
		Results: results,
	}, errs
}

func (m *Matrix) Run(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	errs := sumError(nil)
	entries := make([]*MatrixEntry, 0, len(m.CRIs)*len(m.OCIs))
	for _, cri := range m.CRIs {
		for _, oci := range m.OCIs {
			entry, err := m.createEntry(cri, oci)
			if err != nil {
				errs = append(errs, err)
				break
			}
			entries = append(entries, entry)
		}
	}
	if err := encoder.Encode(entries); err != nil {
		errs = append(errs, err)
	}
	return errs
}

type Suite []Benchmark

func (bs Suite) Name() string {
	names := make([]string, len(bs))
	for i := range bs {
		names[i] = bs[i].Name()
	}
	return fmt.Sprintf("Suite [%s]", strings.Join(names, ", "))
}

func (bs Suite) Run(client *runtime.Client, handler string) (interface{}, error) {
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

// ID computes a unique string for this benchmark.
func ID(b Benchmark) string {
	return fmt.Sprintf("benchmark.%s.%s", b.Name(), runtime.NewUUID())
}
