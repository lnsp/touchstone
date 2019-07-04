package cmd

import (
	"fmt"

	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const sampleLinuxImage = "docker.io/lnsp/sysbench:latest"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the benchmark suite.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := runtime.NewClient(util.GetCRIEndpoint(cri))
		if err != nil {
			logrus.WithError(err).Fatal("failed to create client")
		}
		logrus.Info("pulled sysbench image from Docker Hub")
		sandbox := client.InitLinuxSandbox("sysbench-" + runtime.NewUUID())
		if err := client.PullImage(sampleLinuxImage, nil); err != nil {
			logrus.Fatal(err)
		}
		pod, err := client.StartSandbox(sandbox, handler)
		if err != nil {
			logrus.WithError(err).Fatal("failed to create sandbox")
		}
		logrus.WithField("pod", pod).Info("created sandbox")
		container, err := client.CreateContainer(sandbox, pod, "sysbench", sampleLinuxImage, []string{"sysbench", "--test=cpu", "run"})
		if err != nil {
			logrus.WithError(err).Fatal("failed to create container")
		}
		logrus.WithField("container", container).Info("created container")
		if err := client.StartContainer(container); err != nil {
			logrus.WithError(err).Fatal("failed to start container")
		}
		logs, err := client.WaitForLogs(container)
		if err != nil {
			logrus.WithError(err).Fatal("failed to wait for logs")
		}
		fmt.Println(string(logs))
		if err := client.StopAndRemoveContainer(container); err != nil {
			logrus.WithError(err).Fatal("failed to stop container")
		}
		if err := client.StopAndRemoveSandbox(pod); err != nil {
			logrus.WithError(err).Fatal("failed to stop and remove sandbox")
		}
		logrus.Info("benchmark finished")
	},
}
