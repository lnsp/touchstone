package suites

import (
	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
)

var Performance = []benchmark.Benchmark{
	&MemoryThroughput{},
	&CPUTime{},
}

const defaultSysbenchImage = "lnsp/sysbench:latest"

// RunInSysbench executes a specific sysbench benchmark and returns the application logs.
func RunInSysbench(bm benchmark.Benchmark, client *runtime.Client, handler string, args []string) ([]byte, error) {
	var (
		sandboxID   = benchmark.ID(bm)
		containerID = benchmark.ID(bm)
	)
	// Pull image
	if err := client.PullImage(defaultSysbenchImage, nil); err != nil {
		return nil, err
	}
	// Perform benchmark
	sandbox := client.InitLinuxSandbox(sandboxID)
	pod, err := client.StartSandbox(sandbox, handler)
	if err != nil {
		return nil, err
	}
	container, err := client.CreateContainer(sandbox, pod, containerID, defaultSysbenchImage, args)
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
	// Cleanup container and sandbox
	if err := client.StopAndRemoveContainer(container); err != nil {
		return nil, err
	}
	if err := client.StopAndRemoveSandbox(pod); err != nil {
		return nil, err
	}
	return logs, nil
}

// MemoryThroughput measures the total memory operations and min/avg/max memory latency.
type MemoryThroughput struct{}

func (MemoryThroughput) Name() string {
	return "performance.memory.throughput"
}

func (bm *MemoryThroughput) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	logs, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=memory",
		"--memory-block-size=1M", "--memory-total-size=100G",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"TotalTime":  util.ParsePrefixedLine(logs, "total time:"),
		"MinLatency": util.ParsePrefixedLine(logs, "min:"),
		"AvgLatency": util.ParsePrefixedLine(logs, "avg:"),
		"MaxLatency": util.ParsePrefixedLine(logs, "max:"),
	}, nil
}

// CPUTime measures the total time taken by a CPU heavy task.
type CPUTime struct{}

func (CPUTime) Name() string {
	return "performance.cpu.time"
}

func (bm *CPUTime) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	logs, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test-cpu",
		"run",
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"TotalTime": util.ParsePrefixedLine(logs, "total time:"),
	}, nil
}
