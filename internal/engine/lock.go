//go:build unix
// +build unix

package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const lockFileName = "save.lock"

// Lock holds an open file descriptor used as a process-lifetime lock.
type Lock struct {
	f *os.File
}

// AcquireLock attempts to get an exclusive flock on the lock file.
func AcquireLock() (*Lock, bool, error) {
	dir, err := SaveDir()
	if err != nil {
		return nil, false, err
	}
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, false, err
	}
	lockPath := filepath.Join(dir, lockFileName)
	f, err := os.OpenFile(lockPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, false, err
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		f.Close()
		if err == syscall.EWOULDBLOCK {
			return nil, true, nil
		}
		return nil, false, fmt.Errorf("flock: %w", err)
	}
	f.Truncate(0)
	f.WriteString(fmt.Sprintf("pid=%d\n", os.Getpid()))
	return &Lock{f: f}, false, nil
}

// Release frees the file lock.
func (l *Lock) Release() {
	if l != nil && l.f != nil {
		syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
		l.f.Close()
	}
}

// ReadLockOwner returns the PID written in the lock file.
func ReadLockOwner() int {
	dir, err := SaveDir()
	if err != nil {
		return 0
	}
	data, err := os.ReadFile(filepath.Join(dir, lockFileName))
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "pid=") {
			pid, _ := strconv.Atoi(strings.TrimPrefix(line, "pid="))
			return pid
		}
	}
	return 0
}
