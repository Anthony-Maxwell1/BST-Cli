//go:build !windows

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func Run() error {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)
	corePath := filepath.Join(exeDir, "core", "BST-Core")
	pidFile := filepath.Join(exeDir, "core", "daemon.pid")

	if _, err := os.Stat(pidFile); err == nil {
        fmt.Println("Another instance is already running. Exiting.")
        os.Exit(1)
    }

	cmd := exec.Command(corePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	// Save PID
	pid := cmd.Process.Pid
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	fmt.Printf("Daemon started in background with PID %d\n", pid)
	return nil
}

func Stop() error {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)
	pidFile := filepath.Join(exeDir, "core", "daemon.pid")
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("failed to read PID file: %w", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("invalid PID in file: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	err = process.Signal(syscall.SIGTERM)

	if err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}

	os.Remove(pidFile)
	fmt.Println("Daemon stopped")
	return nil
}
