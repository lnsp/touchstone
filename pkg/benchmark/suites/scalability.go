package suites

import (
	"fmt"
	"time"

	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/runtime"
)

var Scalability = []benchmark.Benchmark{
	&StartupScalability{Scale: 5},
	&StartupScalability{Scale: 10},
	&StartupScalability{Scale: 50},
}

type StartupScalability struct {
	Scale int
}

func (bm *StartupScalability) Name() string {
	return fmt.Sprintf("scalability.runtime.%d", bm.Scale)
}

func (bm *StartupScalability) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	var (
		sandboxNames   = make([]string, bm.Scale)
		containerNames = make([]string, bm.Scale)
		podIDs         = make([]string, bm.Scale)
		containerIDs   = make([]string, bm.Scale)
		image          = "busybox:latest"
	)
	err := client.PullImage(image, nil)
	if err != nil {
		return nil, err
	}
	for i := 0; i < bm.Scale; i++ {
		containerNames[i] = benchmark.ID(bm)
		sandboxNames[i] = benchmark.ID(bm)
	}
	start := time.Now()
	for i := 0; i < bm.Scale; i++ {
		sandbox := client.InitLinuxSandbox(sandboxNames[i])
		podIDs[i], err = client.StartSandbox(sandbox, handler)
		if err != nil {
			return nil, err
		}
		containerIDs[i], err = client.CreateContainer(sandbox, podIDs[i], containerNames[i], image, []string{"sleep", "1000000"})
		if err != nil {
			return nil, err
		}
		if err := client.StartContainer(containerIDs[i]); err != nil {
			return nil, err
		}
	}
	end := time.Now()
	// cleanup
	for i := 0; i < bm.Scale; i++ {
		if err := client.StopAndRemoveContainer(containerIDs[i]); err != nil {
			return nil, err
		}
		if err := client.StopAndRemoveSandbox(podIDs[i]); err != nil {
			return nil, err
		}
	}
	return benchmark.ValueReport{
		"TotalTime": end.Sub(start).Seconds(),
	}, nil
}

func (StartupScalability) Labels() []string {
	return []string{"TotalTime"}
}
