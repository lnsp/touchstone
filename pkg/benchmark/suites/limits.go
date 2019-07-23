package suites

import (
	"time"

	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

var Limits = []benchmark.Benchmark{
	&CPULimits{},
	&CPUScalingLimits{},
}

// RunInSysbench executes a specific sysbench benchmark and returns the application logs.
func RunInSysbenchWithResources(bm benchmark.Benchmark, client *runtime.Client, handler string, args []string, resources *runtimeapi.LinuxContainerResources) ([]byte, error) {
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
	container, err := client.CreateContainerWithResources(sandbox, pod, containerID, defaultSysbenchImage, args, resources)
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

func RunInSysbenchWithScalingResources(bm benchmark.Benchmark, client *runtime.Client, handler string, args []string, resources []*runtimeapi.LinuxContainerResources, interval time.Duration) ([]byte, error) {
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
	container, err := client.CreateContainerWithResources(sandbox, pod, containerID, defaultSysbenchImage, args, resources[0])
	if err != nil {
		return nil, err
	}
	if err := client.StartContainer(container); err != nil {
		return nil, err
	}
	for i := 1; i < len(resources); i++ {
		state, err := client.State(container)
		if err != nil {
			return nil, err
		}
		if state != runtimeapi.ContainerState_CONTAINER_RUNNING {
			break
		}
		if err := client.UpdateContainerResources(container, resources[i]); err != nil {
			return nil, err
		}
		<-time.After(interval)
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

// CPULimits measures the total time taken by a CPU heavy task.
type CPULimits struct{}

func (CPULimits) Name() string {
	return "limits.cpu.time"
}

func (bm *CPULimits) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	logs, err := RunInSysbenchWithResources(bm, client, handler, []string{
		"sysbench", "--test=cpu",
		"--cpu-max-prime=20000",
		"--num-threads=1", "run",
	}, &runtimeapi.LinuxContainerResources{
		CpuPeriod: 100000,
		CpuQuota:  10000,
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"TotalTime": util.ParsePrefixedLine(logs, "total time:"),
	}, nil
}

func (CPULimits) Labels() []string {
	return []string{"TotalTime"}
}

// CPUScalingLimits measures the total time taken by a CPU heavy task.
type CPUScalingLimits struct{}

func (CPUScalingLimits) Name() string {
	return "limits.cpu.scaling"
}

func (bm *CPUScalingLimits) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	resources := make([]*runtimeapi.LinuxContainerResources, 10)
	for i := 0; i < 10; i++ {
		resources[i] = &runtimeapi.LinuxContainerResources{
			CpuPeriod: 100000,
			CpuQuota:  10000 * int64(i+1),
		}
	}
	logs, err := RunInSysbenchWithScalingResources(bm, client, handler, []string{
		"sysbench", "--test=cpu",
		"--cpu-max-prime=5000",
		"--num-threads=1", "run",
	}, resources, time.Second)
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"TotalTime": util.ParsePrefixedLine(logs, "total time:"),
	}, nil
}

func (CPUScalingLimits) Labels() []string {
	return []string{"TotalTime"}
}
