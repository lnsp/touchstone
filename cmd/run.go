package cmd

import (
	"time"

	"github.com/lnsp/touchstone/pkg/framework"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const sampleLinuxImage = "docker.io/library/busybox:latest"

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the benchmark suite.",
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, _ := cmd.Flags().GetString("endpoint")
		client, err := framework.NewClient(endpoint)
		if err != nil {
			logrus.WithError(err).Fatal("failed to create client")
		}
		if err := client.PullImage(sampleLinuxImage, nil); err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("pulled busybox image from Docker Hub")
		sandbox := client.InitLinuxSandbox("busybox-" + framework.NewUUID())
		pod, err := client.StartSandbox(sandbox)
		if err != nil {
			logrus.WithError(err).Fatal("failed to create sandbox")
		}
		logrus.WithField("pod", pod).Info("created sandbox")
		container, err := client.CreateContainer(sandbox, pod, "sandboxed_busybox", sampleLinuxImage, []string{"top"})
		if err != nil {
			logrus.WithError(err).Fatal("failed to create container")
		}
		logrus.WithField("container", container).Info("created container")
		if err := client.StartContainer(container); err != nil {
			logrus.WithError(err).Fatal("failed to start container")
		}
		<-time.After(time.Minute)
		if err := client.StopContainer(container); err != nil {
			logrus.WithError(err).Fatal("failed to stop container")
		}
		if err := client.StopAndRemoveSandbox(pod); err != nil {
			logrus.WithError(err).Fatal("failed to stop and remove sandbox")
		}
		logrus.Info("benchmark finished")
	},
}
