package suites

import (
	"time"

	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/runtime"
)

var Operations = []benchmark.Benchmark{
	&ContainerLifecycle{},
}

type ContainerLifecycle struct{}

func (ContainerLifecycle) Name() string {
	return "operations.container.lifecycle"
}

func (bm *ContainerLifecycle) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	var (
		sandboxID                    = benchmark.ID(bm)
		containerID                  = benchmark.ID(bm)
		image                        = "busybox:latest"
		beginStartup, endStartup     time.Time // measuring total time
		beginSandbox, endSandbox     time.Time // measuring create sandbox and container
		beginContainer, endContainer time.Time // measuring start containerr
		beginShutdown, endShutdown   time.Time // measuring stop container & sandbox
	)
	// Pull image
	if err := client.PullImage(image, nil); err != nil {
		return nil, err
	}
	// Perform benchmark
	sandbox := client.InitLinuxSandbox(sandboxID)
	beginStartup = time.Now()
	beginSandbox = time.Now()
	pod, err := client.StartSandbox(sandbox, handler)
	if err != nil {
		return nil, err
	}
	container, err := client.CreateContainer(sandbox, pod, containerID, image, []string{"sleep", "60"})
	if err != nil {
		return nil, err
	}
	endSandbox = time.Now()
	beginContainer = time.Now()
	if err := client.StartContainer(container); err != nil {
		return nil, err
	}
	endContainer = time.Now()
	endStartup = time.Now()
	beginShutdown = time.Now()
	// Cleanup container and sandbox
	if err := client.StopAndRemoveContainer(container); err != nil {
		return nil, err
	}
	if err := client.StopAndRemoveSandbox(pod); err != nil {
		return nil, err
	}
	endShutdown = time.Now()
	return benchmark.ValueReport{
		"Startup": endStartup.Sub(beginStartup).Seconds(),
		"Create":  endSandbox.Sub(beginSandbox).Seconds(),
		"Start":   endContainer.Sub(beginContainer).Seconds(),
		"Destroy": endShutdown.Sub(beginShutdown).Seconds(),
	}, nil
}
