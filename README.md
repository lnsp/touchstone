# touchstone

Toolset for benchmarking CRI-compatible container runtimes.

## Installation
```bash
$ go get github.com/lnsp/touchstone
```

## Usage
```bash
$ export CONTAINER_RUNTIME_ENDPOINT=unix:///var/run/containerd/containerd.sock
$ touchstone run
```

