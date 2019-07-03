package benchmark

import (
	"time"

	"github.com/lnsp/touchstone/pkg/runtime"
)

// SuiteStartupTime tests the startup times of different container setups.
var SuiteStartupTime = BenchmarkSuite([]Benchmark{
	&totalStartupTimeBenchmark{},
})

type totalStartupTimeBenchmark struct{}

func (totalStartupTimeBenchmark) Name() string {
	return "TotalStartupTime"
}

func (bm *totalStartupTimeBenchmark) Run(client *runtime.Client, handler string) (interface{}, error) {
	report := struct {
		TotalTime int64 `json:"total_time"`
	}{}
	var (
		sandboxID   = getBenchmarkID(bm)
		containerID = getBenchmarkID(bm)
		image       = "busybox:latest"
	)
	// Pull image
	if err := client.PullImage(image, nil); err != nil {
		return nil, err
	}
	// Perform benchmark
	sandbox := client.InitLinuxSandbox(sandboxID)
	start := time.Now()
	pod, err := client.StartSandbox(sandbox, handler)
	if err != nil {
		return nil, err
	}
	container, err := client.CreateContainer(sandbox, pod, containerID, image, []string{"sleep", "60"})
	if err != nil {
		return nil, err
	}
	if err := client.StartContainer(container); err != nil {
		return nil, err
	}
	report.TotalTime = int64(time.Since(start))
	// Cleanup container and sandbox
	if err := client.StopAndRemoveContainer(container); err != nil {
		return nil, err
	}
	if err := client.StopAndRemoveSandbox(pod); err != nil {
		return nil, err
	}
	return report, nil
}
