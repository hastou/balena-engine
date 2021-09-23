package main

import (
	"unsafe"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// Scheduling policies.
const (
	SCHED_OTHER         = 0
	SCHED_FIFO          = 1
	SCHED_RR            = 2
	SCHED_BATCH         = 3
	SCHED_RESET_ON_FORK = 0x40000000 // Meant to be ORed with the others
)

// setScheduler sets the scheduling policy and priority for the current process.
func setScheduler(policy, prio int) {
	type sched_param struct {
		sched_priority int
	}
	s := &sched_param{int(prio)}
	p := unsafe.Pointer(s)
	_, _, errno := unix.Syscall(unix.SYS_SCHED_SETSCHEDULER, uintptr(0), uintptr(policy), uintptr(p))
	if errno != 0 {
		logrus.Errorf("Syscall SYS_SCHED_SETSCHEDULER failed with policy=%v prio=%v: %v\n",
			policy, prio, errno)
	}
}
