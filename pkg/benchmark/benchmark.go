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

type Attempt struct {
	Count int
	Suite Benchmark
}

func (b *Attempt) Name() string {
	return fmt.Sprintf("%s [%dx]", b.Suite.Name(), b.Count)
}

func (b *Attempt) Run(client *runtime.Client, handler string) (interface{}, error) {
	errs := make([]error, 0)
	results := make([]interface{}, 0, b.Count)
	for i := 0; i < b.Count; i++ {
		logrus.WithFields(logrus.Fields{
			"attempt": i,
			"max":     b.Count,
			"suite":   b.Suite.Name(),
		}).Info("running benchmark attempt")
		result, err := b.Suite.Run(client, handler)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results = append(results, result)
	}
	if len(errs) > 0 {
		return nil, aggregatedError(errs)
	}
	return results, nil
}

type aggregatedError []error

func (e aggregatedError) Error() string {
	s := ""
	for _, err := range e {
		s += err.Error() + "\n"
	}
	return s
}

type Matrix struct {
	CRIs            []string
	RuntimeHandlers []string
	Suite           Benchmark
}

func (m *Matrix) Run(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	errs := make([]error, 0)

	for _, cri := range m.CRIs {
		for _, handler := range m.RuntimeHandlers {
			func() {
				logrus.WithFields(logrus.Fields{
					"cri":     cri,
					"handler": handler,
					"suite":   m.Suite.Name(),
				}).Info("running benchmark matrix entry")
				client, err := runtime.NewClient(util.GetCRIEndpoint(cri))
				if err != nil {
					errs = append(errs, fmt.Errorf("(%s / %s): %v", cri, handler, err))
					return
				}
				defer client.Close()
				result, err := m.Suite.Run(client, handler)
				if err != nil {
					errs = append(errs, fmt.Errorf("(%s / %s): %v", cri, handler, err))
					return
				}
				wrapper := struct {
					CRI            string `json:"cri_endpoint"`
					RuntimeHandler string `json:"runtime_handler"`
					Result         interface{}
				}{
					CRI:            cri,
					RuntimeHandler: handler,
					Result:         result,
				}
				if err := encoder.Encode(wrapper); err != nil {
					errs = append(errs, fmt.Errorf("(%s / %s): %v", cri, handler, err))
				}
			}()
		}
	}
	if len(errs) > 0 {
		return aggregatedError(errs)
	}
	return nil
}

type Benchmark interface {
	Name() string
	Run(client *runtime.Client, handler string) (interface{}, error)
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

// getBenchmarkID computes a unique string for this benchmark.
func getBenchmarkID(b Benchmark) string {
	return fmt.Sprintf("Benchmark-%s-%s", b.Name(), runtime.NewUUID())
}
