package framework

import (
	internalapi "k8s.io/cri-api"
	"k8s.io/kubernetes/pkg/kubelet/remote"
)

type APIClient struct {
	Runtime *internalapi.RuntimeServiceClient
	Image   *internalapi.ImageServiceClient
}

func NewClient(addr string) *APIClient {
	runtimeClient := &APIClient{

	}
	return runtimeClient
}
