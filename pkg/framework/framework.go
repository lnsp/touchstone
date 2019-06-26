package framework

import (
	"time"

	"github.com/pkg/errors"
	internalapi "k8s.io/cri-api/pkg/apis"
	"k8s.io/kubernetes/pkg/kubelet/remote"
)

type APIClient struct {
	Runtime internalapi.RuntimeService
	Image   internalapi.ImageManagerService
}

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
