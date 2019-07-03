package util

import "testing"

func TestGetCRIEndpoint(t *testing.T) {
	tt := []struct {
		Name string
		Path string
	}{
		{"crio", "unix:///var/run/crio/crio.sock"},
		{"containerd", "unix:///var/run/containerd/containerd.sock"},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			path := GetCRIEndpoint(tc.Name)
			if path != tc.Path {
				t.Errorf("expected %s, got %s", tc.Path, path)
			}
		})
	}
}
