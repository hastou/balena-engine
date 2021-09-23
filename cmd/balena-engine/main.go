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
		fmt.Printf("Process %v/%v is balenad\n", os.Getpid(), os.Getppid())
		policy, prio := schedArgsFromEnv("DAEMON", SCHED_RR, 50)
		if policy >= 0 {
			setScheduler(policy, prio)
		}
		dockerd.Main()
	case "balena-containerd", "balena-engine-containerd":
		fmt.Printf("Process %v/%v is balena-containerd\n", os.Getpid(), os.Getppid())
		policy, prio := schedArgsFromEnv("CONTAINERD", SCHED_RR, 50)
		if policy >= 0 {
			setScheduler(policy|SCHED_RESET_ON_FORK, prio)
		}
		containerd.Main()
	case "balena-containerd-shim", "balena-engine-containerd-shim":
		containerdShim.Main()
	case "balena-containerd-ctr", "balena-engine-containerd-ctr":
		ctr.Main()
	case "balena-runc", "balena-engine-runc":
		runc.Main()
	case "balena-proxy", "balena-engine-proxy":
		proxy.Main()
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command: %v\n", command)
		os.Exit(1)
	}
}
