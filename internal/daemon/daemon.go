package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const pidFile = "core/daemon.pid"

func Run() error {
	binary := "./core/BST-Core"

	if runtime.GOOS == "windows" {
		binary = "./core/BST-Core.exe"
	}

	cmd := exec.Command(binary)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	pid := cmd.Process.Pid
	os.WriteFile(pidFile, []byte(fmt.Sprint(pid)), 0644)

	fmt.Println("Daemon started with PID", pid)
	return nil
}

func Stop() error {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return err
	}

	var pid int
	fmt.Sscanf(string(data), "%d", &pid)

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Kill()
	if err != nil {
		return err
	}

	os.Remove(pidFile)
	fmt.Println("Daemon stopped")
	return nil
}
	