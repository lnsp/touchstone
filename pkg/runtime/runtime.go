package runtime

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	internalapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/remote"
)

const maxRemovalAttempts = 10
const maxRemovalTimeout = 10

var uuidLock sync.Mutex
var lastUUID uuid.UUID

// NewUUID actively waits for a new UUID and returns the uuid string.
func NewUUID() string {
	uuidLock.Lock()
	defer uuidLock.Unlock()

	result := uuid.NewUUID()
	for uuid.Equal(lastUUID, result) == true {
		result = uuid.NewUUID()
	}
	lastUUID = result
	return result.String()
}

const defaultNamespace = "touchstone"

// Client is an implementation of a CRI API client.
type Client struct {
	Runtime internalapi.RuntimeService
	Image   internalapi.ImageManagerService
}

var defaultLinuxPodLabels = map[string]string{}

// CreateContainer runs a container image. It returns the container ID.
func (api *Client) CreateContainer(sandbox *runtimeapi.PodSandboxConfig, pod, name, image string, command []string) (string, error) {
	container := &runtimeapi.ContainerConfig{
		Metadata: &runtimeapi.ContainerMetadata{
			Name:    name,
			Attempt: 0,
		},
		Image: &runtimeapi.ImageSpec{
			Image: image,
		},
		LogPath: "/var/log/" + pod + "_" + name + ".log",
		Linux:   &runtimeapi.LinuxContainerConfig{},
	}
	if command != nil {
		container.Command = command
	}
	return api.Runtime.CreateContainer(pod, container, sandbox)
}

// WaitForLogs waits for the container to exit and returns the logs as a slice of bytes.
func (api *Client) WaitForLogs(container string) ([]byte, error) {
	for {
		status, err := api.Status(container)
		if err != nil {
			return nil, err
		}
		if status.State >= 2 {
			break
		}
		<-time.After(time.Second)
	}
	buf := &bytes.Buffer{}
	if err := api.Logs(container, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Logs fetches the logs of the container.
func (api *Client) Logs(container string, writer io.Writer) error {
	status, err := api.Status(container)
	if err != nil {
		return err
	}
	logPath := status.GetLogPath()
	if logPath == "" {
		return errors.New("missing log path")
	}

	f, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %v", logPath, err)
	}
	defer f.Close()

	// Start parsing the logs.
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		items := bytes.SplitN(line, []byte(" "), 4)
		if len(items) == 4 {
			fmt.Fprintln(writer, string(items[3]))
		}
	}
	return nil
}

// Status fetches the status of a container.
func (api *Client) Status(container string) (*runtimeapi.ContainerStatus, error) {
	return api.Runtime.ContainerStatus(container)
}

// StartContainer starts a new container instance.
func (api *Client) StartContainer(container string) error {
	return api.Runtime.StartContainer(container)
}

// StopContainer stops the container instance.
func (api *Client) StopContainer(container string) error {
	return api.Runtime.StopContainer(container, maxRemovalTimeout)
}

// StartSandbox starts up the pod sandbox. It returns the pod sandbox ID.
func (api *Client) StartSandbox(sandbox *runtimeapi.PodSandboxConfig, runtime string) (string, error) {
	return api.Runtime.RunPodSandbox(sandbox, runtime)
}

// StopAndRemoveContainer stops and removes a container.
func (api *Client) StopAndRemoveContainer(container string) (err error) {
	for attempt := 0; attempt < maxRemovalAttempts; attempt++ {
		err = api.Runtime.StopContainer(container, maxRemovalTimeout)
		if err != nil {
			continue
		}
		err = api.Runtime.RemoveContainer(container)
		if err != nil {
			continue
		}
		return nil
	}
	return errors.Errorf("stop-remove container failed: %v", err)
}

// StopAndRemoveSandbox stops and removes the given pod sandbox.
func (api *Client) StopAndRemoveSandbox(pod string) (err error) {
	for attempt := 0; attempt < maxRemovalAttempts; attempt++ {
		err = api.Runtime.StopPodSandbox(pod)
		if err != nil {
			continue
		}
		err = api.Runtime.RemovePodSandbox(pod)
		if err != nil {
			continue
		}
		return nil
	}
	return errors.Errorf("stop-remove pod failed: %v", err)
}

// InitLinuxSandbox creates a new pod sandbox configuration.
func (api *Client) InitLinuxSandbox(name string) *runtimeapi.PodSandboxConfig {
	return &runtimeapi.PodSandboxConfig{
		Metadata: &runtimeapi.PodSandboxMetadata{
			Name:      name,
			Uid:       NewUUID(),
			Namespace: defaultNamespace,
			Attempt:   1,
		},
		Linux:  &runtimeapi.LinuxPodSandboxConfig{},
		Labels: defaultLinuxPodLabels,
	}
}

// PullImage instructs the CRI to pull an image from a public repository.
func (api *Client) PullImage(image string, sandbox *runtimeapi.PodSandboxConfig) error {
	if !strings.Contains(image, ":") {
		image = image + ":latest"
	}
	imageSpec := &runtimeapi.ImageSpec{
		Image: image,
	}
	_, err := api.Image.PullImage(imageSpec, nil, sandbox)
	if err != nil {
		return err
	}
	return nil
}

func (api *Client) Close() {
	// TODO: close TCP connections
}

// NewClient instantiates a new API client.
func NewClient(addr string) (*Client, error) {
	logrus.WithField("addr", addr).Info("connecting to CRI endpoint")
	runtimeSvc, err := remote.NewRemoteRuntimeService(addr, 10*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}
	imageSvc, err := remote.NewRemoteImageService(addr, 10*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}
	runtimeClient := &Client{
		Runtime: runtimeSvc,
		Image:   imageSvc,
	}
	return runtimeClient, nil
}
