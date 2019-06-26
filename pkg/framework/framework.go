package framework

import (
	"strings"
	"sync"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	internalapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/remote"
)

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

// APIClient is an implementation of a CRI API client.
type APIClient struct {
	Runtime internalapi.RuntimeService
	Image   internalapi.ImageManagerService
}

var defaultLinuxPodLabels = map[string]string{}

// CreateContainer runs a container image. It returns the container ID.
func (api *APIClient) CreateContainer(sandbox *runtimeapi.PodSandboxConfig, pod, name, image string, command []string) (string, error) {
	container := &runtimeapi.ContainerConfig{
		Metadata: &runtimeapi.ContainerMetadata{
			Name:    name,
			Attempt: 0,
		},
		Image: &runtimeapi.ImageSpec{
			Image: image,
		},
		Linux: &runtimeapi.LinuxContainerConfig{},
	}
	if command != nil {
		container.Command = command
	}
	return api.Runtime.CreateContainer(pod, container, sandbox)
}

// StartContainer starts a new container instance.
func (api *APIClient) StartContainer(container string) error {
	return api.Runtime.StartContainer(container)
}

// StopContainer stops the container instance.
func (api *APIClient) StopContainer(container string) error {
	return api.Runtime.StopContainer(container, 60)
}

// StartSandbox starts up the pod sandbox. It returns the pod sandbox ID.
func (api *APIClient) StartSandbox(sandbox *runtimeapi.PodSandboxConfig) (string, error) {
	return api.Runtime.RunPodSandbox(sandbox, "")
}

// StopAndRemoveSandbox stops and removes the given pod sandbox.
func (api *APIClient) StopAndRemoveSandbox(pod string) error {
	err := api.Runtime.StopPodSandbox(pod)
	if err != nil {
		return err
	}
	return api.Runtime.RemovePodSandbox(pod)
}

// InitLinuxSandbox creates a new pod sandbox configuration.
func (api *APIClient) InitLinuxSandbox(name string) *runtimeapi.PodSandboxConfig {
	return &runtimeapi.PodSandboxConfig{
		Metadata: &runtimeapi.PodSandboxMetadata{
			Name:      name,
			Uid:       NewUUID(),
			Namespace: defaultNamespace,
			Attempt:   0,
		},
		Linux:  &runtimeapi.LinuxPodSandboxConfig{},
		Labels: defaultLinuxPodLabels,
	}
}

// PullImage instructs the CRI to pull an image from a public repository.
func (api *APIClient) PullImage(image string, sandbox *runtimeapi.PodSandboxConfig) error {
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

// NewClient instantiates a new API client.
func NewClient(addr string) (*APIClient, error) {
	runtimeSvc, err := remote.NewRemoteRuntimeService(addr, time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}
	imageSvc, err := remote.NewRemoteImageService(addr, time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}
	runtimeClient := &APIClient{
		Runtime: runtimeSvc,
		Image:   imageSvc,
	}
	return runtimeClient, nil
}
