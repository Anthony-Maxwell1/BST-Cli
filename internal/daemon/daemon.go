package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
)

const pidFile = "core/daemon.pid"

func Run() error {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)
	corePath := filepath.Join(exeDir, "core", "BST-Core.exe")	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	cmd := exec.Command(corePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if runtime.GOOS == "windows" {
		// Windows detach constants
		const (
			CREATE_NEW_PROCESS_GROUP = 0x00000200
			DETACHED_PROCESS         = 0x00000008
		)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: CREATE_NEW_PROCESS_GROUP | DETACHED_PROCESS,
		}
	} else {
		// Unix: leave SysProcAttr nil
		cmd.SysProcAttr = nil
	}

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

	if runtime.GOOS == "windows" {
		err = process.Kill()
	} else {
		err = process.Signal(syscall.SIGTERM)
	}

	if err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}

	os.Remove(pidFile)
	fmt.Println("Daemon stopped")
	return nil
}
