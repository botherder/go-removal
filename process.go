package removal

import (
	"time"
	"errors"
	"github.com/shirou/gopsutil/process"
)

func FindProcessByName(name string) ([]int32, error) {
	var matches []int32

	// We loop over all processes.
	pids, err := process.Pids()
	if err != nil {
		return matches, err
	}

	for _, pid := range(pids) {
		// Open process handle.
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		// Get process name.
		procName, err := proc.Name()
		if err != nil {
			continue
		}

		// We check if the current process is the one we're looking for.
		if procName == name {
			matches = append(matches, pid)
		}
	}

	return matches, nil
}

func GetProcessExe(pid int32) (string, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return "", err
	}

	return proc.Exe()
}

func KillProcessByPid(pid int32) error {
	// Open process handle.
	proc, err := process.NewProcess(pid)
	if err != nil {
		return err
	}

	// Terminate the process.
	err = proc.Terminate()
	if err != nil {
		return err
	}

	return nil
}

func WaitForProcessToDie(pid int32) error {
	timeout := time.After(30 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return errors.New("timed out")
		case <-tick:
			exists, _ := process.PidExists(pid)
			if exists == false {
				return nil
			}
		}
	}
}
