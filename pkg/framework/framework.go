package framework

type CRIClient struct {
	runtime *runtimeapi.RuntimeServiceClient
	image   *imageapi.ImageServiceClient
}
