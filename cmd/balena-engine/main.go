package main

import (
	"fmt"

	"github.com/containerd/containerd/cmd/containerd"
	containerdShim "github.com/containerd/containerd/cmd/containerd-shim"
	"github.com/containerd/containerd/cmd/ctr"
	"github.com/docker/cli/cmd/docker"
	"github.com/docker/docker/cmd/dockerd"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/libnetwork/cmd/proxy"
	"github.com/opencontainers/runc"

	"os"
	filepath "path/filepath"
)

func main() {
	if reexec.Init() {
		return
	}

	command := filepath.Base(os.Args[0])

	switch command {
	case "balena", "balena-engine":
		docker.Main()
	case "balenad", "balena-engine-daemon":
		setScheduler(SCHED_RR, 35)
		dockerd.Main()
	case "balena-containerd", "balena-engine-containerd":
		setScheduler(SCHED_RR, 35)
		containerd.Main()
	case "balena-containerd-shim", "balena-engine-containerd-shim":
		containerdShim.Main()
	case "balena-containerd-ctr", "balena-engine-containerd-ctr":
		ctr.Main()
	case "balena-runc", "balena-engine-runc":
		setScheduler(SCHED_OTHER, 35)
		runc.Main()
	case "balena-proxy", "balena-engine-proxy":
		proxy.Main()
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command: %v\n", command)
		os.Exit(1)
	}
}
