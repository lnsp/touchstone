package benchmark

import (
	"fmt"

	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
)

var SuitePerformance = Suite([]Benchmark{
	&totalSysbenchCPUBenchmark{},
})

type totalSysbenchCPUBenchmark struct{}

func (totalSysbenchCPUBenchmark) Name() string {
	return "TotalSysbenchCPUTime"
}

func (bm *totalSysbenchCPUBenchmark) Run(client *runtime.Client, handler string) (interface{}, error) {
	report := struct {
		TotalTime string `json:"total_time"`
	}{}
	var (
		sandboxID   = getBenchmarkID(bm)
		containerID = getBenchmarkID(bm)
		image       = "lnsp/sysbench:latest"
	)
	// Pull image
	if err := client.PullImage(image, nil); err != nil {
		return nil, err
	}
	// Perform benchmark
	sandbox := client.InitLinuxSandbox(sandboxID)
	pod, err := client.StartSandbox(sandbox, handler)
	if err != nil {
		return nil, err
	}
	container, err := client.CreateContainer(sandbox, pod, containerID, image, []string{"sysbench", "--test=cpu", "run"})
	if err != nil {
		return nil, err
	}
	if err := client.StartContainer(container); err != nil {
		return nil, err
	}
	logs, err := client.WaitForLogs(container)
	if err != nil {
		return nil, err
	}
	// search for 'total time:'
	totalTime, err := util.FindPrefixedLine(logs, "total time:")
	if err != nil {
		return nil, fmt.Errorf("failed to parse logs: %v", err)
	}
	report.TotalTime = totalTime
	// Cleanup container and sandbox
	if err := client.StopAndRemoveContainer(container); err != nil {
		return nil, err
	}
	if err := client.StopAndRemoveSandbox(pod); err != nil {
		return nil, err
	}
	return report, nil
}
