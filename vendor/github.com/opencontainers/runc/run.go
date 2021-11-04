//go:build linux
// +build linux

package runc

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/urfave/cli"
	"golang.org/x/sys/unix"
)

///////////////////////////////////////

// Scheduling policies.
const (
	SCHED_NORMAL        = 0
	SCHED_FIFO          = 1
	SCHED_RR            = 2
	SCHED_BATCH         = 3
	SCHED_RESET_ON_FORK = 0x40000000 // Meant to be ORed with the others
)

///////////////////////////////////////

// default action is to start a container
var runCommand = cli.Command{
	Name:  "run",
	Usage: "create and run a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The run command creates an instance of a container for a bundle. The bundle
is a directory with a specification file named "` + specConfig + `" and a root
filesystem.

The specification file includes an args parameter. The args parameter is used
to specify command(s) that get run when the container is started. To change the
command(s) that get executed on start, edit the args parameter of the spec. See
"runc spec --help" for more explanation.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: `path to the root of the bundle directory, defaults to the current directory`,
		},
		cli.StringFlag{
			Name:  "console-socket",
			Value: "",
			Usage: "path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal",
		},
		cli.BoolFlag{
			Name:  "detach, d",
			Usage: "detach from the container's process",
		},
		cli.StringFlag{
			Name:  "pid-file",
			Value: "",
			Usage: "specify the file to write the process id to",
		},
		cli.BoolFlag{
			Name:  "no-subreaper",
			Usage: "disable the use of the subreaper used to reap reparented processes",
		},
		cli.BoolFlag{
			Name:  "no-pivot",
			Usage: "do not use pivot root to jail process inside rootfs.  This should be used whenever the rootfs is on top of a ramdisk",
		},
		cli.BoolFlag{
			Name:  "no-new-keyring",
			Usage: "do not create a new session keyring for the container.  This will cause the container to inherit the calling processes session key",
		},
		cli.IntFlag{
			Name:  "preserve-fds",
			Usage: "Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)",
		},
		cli.BoolFlag{
			Name:  "keep-rt-scheduling",
			Usage: "keep the runc process running with a realtime scheduling policy",
		},
	},
	Action: func(context *cli.Context) error {
		if err := checkArgs(context, 1, exactArgs); err != nil {
			return err
		}
		if !context.Bool("keep-rt-scheduling") {
			type sched_param struct {
				sched_priority int
			}
			s := &sched_param{0}
			p := unsafe.Pointer(s)
			_, _, errno := unix.Syscall(unix.SYS_SCHED_SETSCHEDULER, 0, SCHED_NORMAL, uintptr(p))
			if errno != 0 {
				fmt.Fprintf(os.Stderr, "Syscall SYS_SCHED_SETSCHEDULER failed: %v\n", errno)
				// TODO: Return an error or keep running anyway?
			}
		}
		if err := revisePidFile(context); err != nil {
			return err
		}
		spec, err := setupSpec(context)
		if err != nil {
			return err
		}
		status, err := startContainer(context, spec, CT_ACT_RUN, nil)
		if err == nil {
			// exit with the container's exit status so any external supervisor is
			// notified of the exit with the correct exit status.
			os.Exit(status)
		}
		return err
	},
}
