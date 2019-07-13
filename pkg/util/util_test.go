package util

import (
	"testing"
)

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

func TestFindPrefixedLine(t *testing.T) {
	sysbenchOutput := []byte(`sysbench 0.4.12:  multi-threaded system evaluation benchmark

Running the test with following options:
Number of threads: 1

Doing CPU performance benchmark

Threads started!
Done.

Maximum prime number checked in CPU test: 10000


Test execution summary:
    total time:                          10.0634s
    total number of events:              10000
    total time taken by event execution: 10.0610
    per-request statistics:
         min:                                  0.86ms
         avg:                                  1.01ms
         max:                                  2.91ms
         approx.  95 percentile:               1.33ms

Threads fairness:
    events (avg/stddev):           10000.0000/0.00
    execution time (avg/stddev):   10.0610/0.00
`)
	expected := "10.0634s"
	prefix := "total time:"
	value := FindPrefixedLine(sysbenchOutput, prefix)
	if value != expected {
		t.Errorf("expected %s, got %s", expected, value)
	}
}
