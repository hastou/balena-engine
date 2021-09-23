package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Scheduling policies.
const (
	SCHED_NORMAL        = 0
	SCHED_FIFO          = 1
	SCHED_RR            = 2
	SCHED_BATCH         = 3
	SCHED_RESET_ON_FORK = 0x40000000 // Meant to be ORed with the others
)

// setScheduler sets the scheduling policy and priority of all threads in the
// current process.
//
// This function will actually start a goroutine that will periodically reassign
// the requested policy and priority, because the Go runtime can (and does)
// create new threads after the program startup.
func setScheduler(policy, prio int) {
	type sched_param struct {
		sched_priority int
	}
	s := &sched_param{int(prio)}
	p := unsafe.Pointer(s)

	go func() {
		sleepTime := 2 * time.Second
		for {
			tids, err := getAllTasks()
			// TODO: Can we improve error handling? Maybe better logging?
			if err != nil {
				fmt.Fprintf(os.Stderr, "Setting scheduling policy: %v\n", err)
			}
			for _, tid := range tids {
				_, _, errno := unix.Syscall(unix.SYS_SCHED_SETSCHEDULER, uintptr(tid), uintptr(policy), uintptr(p))
				if errno != 0 {
					fmt.Fprintf(os.Stderr, "Syscall SYS_SCHED_SETSCHEDULER failed: %v\n", errno)
				}
			}

			// At least for our usage patterns, the Go runtime tends to create
			// most threads close to the moment the program started. We thus
			// gradually increase the sleep interval.
			fmt.Printf("Process %v/%v sleeping for %v\n", os.Getpid(), os.Getppid(), sleepTime)
			time.Sleep(sleepTime)
			sleepTime *= 2
			if sleepTime > 15*time.Minute {
				sleepTime = 15 * time.Minute
			}
		}
	}()
}

// getAllTasks returns a list with the IDs of all tasks (threads) of the
// currently running process.
func getAllTasks() ([]int, error) {
	files, err := ioutil.ReadDir("/proc/self/task/")
	if err != nil {
		return nil, err
	}

	tids := []int{}
	for _, file := range files {
		tid, err := strconv.Atoi(file.Name())
		if err != nil {
			return tids, err
		}
		tids = append(tids, tid)
	}

	return tids, nil
}

func schedArgsFromEnv(envVar string, defaultPolicy, defaultPrio int) (policy, prio int) {
	policyStr := os.Getenv("BALENA_ENGINE_" + envVar + "_POLICY")
	fmt.Printf("Process %v/%v, policyStr = '%v'\n", os.Getpid(), os.Getppid(), policyStr)
	switch policyStr {
	case "KEEPIT":
		policy = -1
	case "FIFO":
		policy = SCHED_FIFO
	case "RR":
		policy = SCHED_RR
	case "NORMAL":
		policy = SCHED_NORMAL
	default:
		policy = defaultPolicy
	}
	fmt.Printf("Process %v/%v, policy = %v\n", os.Getpid(), os.Getppid(), policy)

	prio = defaultPrio
	prioStr := os.Getenv("BALENA_ENGINE_" + envVar + "_PRIO")
	fmt.Printf("Process %v/%v, prioStr = '%v'\n", os.Getpid(), os.Getppid(), prioStr)
	if prioStr != "" {
		var err error
		prio, err = strconv.Atoi(prioStr)
		if err != nil {
			prio = defaultPrio
		}
	}
	fmt.Printf("Process %v/%v, prio = '%v'\n", os.Getpid(), os.Getppid(), prio)
	return
}
