# touchstone

Toolset for benchmarking CRI-compatible container runtimes.

## Installation
```bash
$ go get github.com/lnsp/touchstone
```

You need to have CRI-O and containerd as well as runc and gVisor set up and ready for running containers. This includes networking configuration.

## Usage
> Note: Please remember that all commands must be run as a privileged user.

```bash
# check connection to crio, containerd
$ touchstone version
touchstone dev-723c8f8a
containerd v1.2.0-621-g04e7747e
cri-o 1.15.1-dev
# run all benchmarks and spill out results in tmp
$ touchstone benchmark -f="suites/*.yaml" -d /tmp/
```

