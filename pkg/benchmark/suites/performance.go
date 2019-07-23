package suites

import (
	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
)

var Performance = []benchmark.Benchmark{
	&MemoryTime{},
	&MemoryMinAvgLatency{},
	&MemoryMaxLatency{},
	&CPUTime{},
	&DiskRead{},
	&DiskWrite{},
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

// DiskWrite measures the total read/write speed.
type DiskWrite struct{}

func (DiskWrite) Name() string {
	return "performance.disk.write"
}

func (bm *DiskWrite) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	seqwr, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=fileio",
		"--file-test-mode=seqwr",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	seqrewr, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=fileio",
		"--file-test-mode=seqrewr",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	rndwr, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=fileio",
		"--file-test-mode=rndwr",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"SeqWrite":   util.ParsePrefixedLine(seqwr, "total time:"),
		"SeqRewrite": util.ParsePrefixedLine(seqrewr, "total time:"),
		"RndWrite":   util.ParsePrefixedLine(rndwr, "total time:"),
	}, nil
}

func (DiskWrite) Labels() []string {
	return []string{"SeqWrite", "SeqRewrite", "RndWrite"}
}

// DiskRead measures the total read/write speed.
type DiskRead struct{}

func (DiskRead) Name() string {
	return "performance.disk.read"
}

func (bm *DiskRead) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	seqrd, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=fileio",
		"--file-test-mode=seqrd",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	rndrd, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=fileio",
		"--file-test-mode=rndrd",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"SeqRead": util.ParsePrefixedLine(seqrd, "total time:"),
		"RndRead": util.ParsePrefixedLine(rndrd, "total time:"),
	}, nil
}

func (DiskRead) Labels() []string {
	return []string{"SeqRead", "RndRead"}
}

// CPUTime measures the total time taken by a CPU heavy task.
type CPUTime struct{}

func (CPUTime) Name() string {
	return "performance.cpu.time"
}

func (bm *CPUTime) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	logs, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=cpu",
		"--cpu-max-prime=20000",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"TotalTime": util.ParsePrefixedLine(logs, "total time:"),
	}, nil
}

func (CPUTime) Labels() []string {
	return []string{"TotalTime"}
}

// MemoryTime measures the total memory operation time.
type MemoryTime struct{}

func (MemoryTime) Name() string {
	return "performance.memory.total"
}

func (bm *MemoryTime) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	logs, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=memory",
		"--memory-block-size=1M", "--memory-total-size=100G",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"TotalTime": util.ParsePrefixedLine(logs, "total time:"),
	}, nil
}

func (MemoryTime) Labels() []string {
	return []string{"TotalTime"}
}

// MemoryMinAvgLatency measures the total memory operations and min/avg/max memory latency.
type MemoryMinAvgLatency struct{}

func (MemoryMinAvgLatency) Name() string {
	return "performance.memory.minavglatency"
}

func (bm *MemoryMinAvgLatency) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	logs, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=memory",
		"--memory-block-size=1M", "--memory-total-size=1G",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"MinLatency": util.ParsePrefixedLine(logs, "min:"),
		"AvgLatency": util.ParsePrefixedLine(logs, "avg:"),
	}, nil
}

func (MemoryMinAvgLatency) Labels() []string {
	return []string{"MinLatency", "AvgLatency"}
}

// MemoryMaxLatency measures the total memory operation time.
type MemoryMaxLatency struct{}

func (MemoryMaxLatency) Name() string {
	return "performance.memory.maxlatency"
}

func (bm *MemoryMaxLatency) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	logs, err := RunInSysbench(bm, client, handler, []string{
		"sysbench", "--test=memory",
		"--memory-block-size=1M", "--memory-total-size=1G",
		"--num-threads=1", "run",
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"MaxLatency": util.ParsePrefixedLine(logs, "max:"),
	}, nil
}

func (MemoryMaxLatency) Labels() []string {
	return []string{"MaxLatency"}
}
