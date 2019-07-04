package benchmark

import (
	"fmt"

	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
)

var SuitePerformance = Suite([]Benchmark{
	&totalSysbenchCPUBenchmark{},
	&memoryThroughputBenchmark{},
})

type memoryThroughputBenchmark struct{}

func (memoryThroughputBenchmark) Name() string {
	return "MemoryThroughput"
}

func (bm *memoryThroughputBenchmark) Run(client *runtime.Client, handler string) (interface{}, error) {
	report := struct {
		TotalOperations string `json:"total_operations"`
		MinLatency      string `json:"min_latency"`
		MaxLatency      string `json:"max_latency"`
		AvgLatency      string `json:"avg_latency"`
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
	container, err := client.CreateContainer(sandbox, pod, containerID, image, []string{"sysbench", "--test=memory", "--memory-block-size=1M", "--memory-total-size=100G", "--num-threads=1", "run"})
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
	report.TotalOperations, err = util.FindPrefixedLine(logs, "total operations:")
	if err != nil {
		return nil, fmt.Errorf("failed to parse logs: %v", err)
	}
	report.MinLatency, err = util.FindPrefixedLine(logs, "min:")
	if err != nil {
		return nil, fmt.Errorf("failed to parse logs: %v", err)
	}
	report.AvgLatency, err = util.FindPrefixedLine(logs, "avg:")
	if err != nil {
		return nil, fmt.Errorf("failed to parse logs: %v", err)
	}
	report.MaxLatency, err = util.FindPrefixedLine(logs, "max:")
	if err != nil {
		return nil, fmt.Errorf("failed to parse logs: %v", err)
	}
	// Cleanup container and sandbox
	if err := client.StopAndRemoveContainer(container); err != nil {
		return nil, err
	}
	if err := client.StopAndRemoveSandbox(pod); err != nil {
		return nil, err
	}
	return report, nil

}

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
