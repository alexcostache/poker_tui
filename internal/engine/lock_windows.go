//go:build windows
// +build windows

package engine

import (
	"fmt"
	"os"
)

const lockFileName = "save.lock"

type Lock struct {
	f *os.File
}

// AcquireLock is a stub on Windows: locking is not supported.
func AcquireLock() (*Lock, bool, error) {
	dir, err := SaveDir()
	if err != nil {
		return nil, false, err
	}
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, false, err
	}
	lockPath := dir + string(os.PathSeparator) + lockFileName
	f, err := os.OpenFile(lockPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, false, err
	}
	// No locking on Windows, just write PID
	f.Truncate(0)
	f.WriteString(fmt.Sprintf("pid=%d\n", os.Getpid()))
	return &Lock{f: f}, false, nil
}

func (l *Lock) Release() {
	if l != nil && l.f != nil {
		l.f.Close()
	}
}

func ReadLockOwner() int {
	// Optionally implement reading PID from file, or just return 0
	return 0
}
