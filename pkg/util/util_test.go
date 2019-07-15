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
	tt := []struct {
		Name   string
		Output []byte
		Filter string
		Value  string
	}{
		{
			Name: "cpu",
			Output: []byte(`sysbench 0.4.12:  multi-threaded system evaluation benchmark

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
`),
			Value:  "10.0634s",
			Filter: "total time:",
		},
		{
			Name: "memory",
			Output: []byte(`sysbench 0.4.12:  multi-threaded system evaluation benchmark

	Running the test with following options:
	Number of threads: 1
	
	Doing memory operations speed test
	Memory block size: 1024K
	
	Memory transfer size: 102400M
	
	Memory operations type: write
	Memory scope type: global
	Threads started!
	Done.
	
	Operations performed: 102400 (10221.21 ops/sec)
	
	102400.00 MB transferred (10221.21 MB/sec)
	
	
	Test execution summary:
		total time:                          10.0184s
		total number of events:              102400
		total time taken by event execution: 10.0103
		per-request statistics:
			 min:                                  0.09ms
			 avg:                                  0.10ms
			 max:                                  1.53ms
			 approx.  95 percentile:               0.13ms
	
	Threads fairness:
		events (avg/stddev):           102400.0000/0.00
		execution time (avg/stddev):   10.0103/0.00
`),
			Filter: "total time:",
			Value:  "10.0184s",
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			result := FindPrefixedLine(tc.Output, tc.Filter)
			if result != tc.Value {
				t.Errorf("expected %s, got %s", tc.Value, result)
			}
		})
	}
}
