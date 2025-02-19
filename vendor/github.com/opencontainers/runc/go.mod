module github.com/opencontainers/runc

go 1.14

require (
	github.com/checkpoint-restore/go-criu/v4 v4.1.0
	github.com/cilium/ebpf v0.0.0-20200702112145-1c8d4c9ef775
	github.com/containerd/console v1.0.0
	github.com/coreos/go-systemd/v22 v22.1.0
	github.com/cyphar/filepath-securejoin v0.2.2
	github.com/docker/go-units v0.4.0
	github.com/godbus/dbus/v5 v5.0.3
	github.com/golang/protobuf v1.4.2
	github.com/moby/sys/mountinfo v0.3.1
	github.com/mrunalp/fileutils v0.0.0-20200520151820-abd8a0e76976
	github.com/opencontainers/runtime-spec v1.0.3-0.20200728170252-4d89ac9fbff6
	github.com/opencontainers/selinux v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/seccomp/libseccomp-golang v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2
	// NOTE: urfave/cli must be <= v1.22.1 due to a regression: https://github.com/urfave/cli/issues/1092
	github.com/urfave/cli v1.22.1
	github.com/vishvananda/netlink v1.1.0
	golang.org/x/sys v0.0.0-20200909081042-eff7692f9009
)
