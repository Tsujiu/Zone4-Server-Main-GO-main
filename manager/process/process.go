package process

import (
	"errors"
	"os"
	"os/exec"
	"sync"
)

var (
	mu      sync.Mutex
	procs   = map[string]*exec.Cmd{}
	running = map[string]bool{}
)

func StartProcess(id, cmdLine string) error {
	mu.Lock()
	defer mu.Unlock()

	if running[id] {
		return errors.New("already running")
	}

	cmd := exec.Command(shell(), shellArg(), cmdLine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	procs[id] = cmd
	running[id] = true

	go func(id string, c *exec.Cmd) {
		_ = c.Wait()
		mu.Lock()
		defer mu.Unlock()
		delete(procs, id)
		running[id] = false
	}(id, cmd)

	return nil
}

func StopProcess(id string) error {
	mu.Lock()
	defer mu.Unlock()

	cmd, ok := procs[id]
	if !ok {
		return errors.New("not running")
	}
	if cmd.Process == nil {
		return errors.New("no process")
	}
	err := cmd.Process.Kill()
	delete(procs, id)
	running[id] = false
	return err
}

func StopAll() {
	mu.Lock()
	defer mu.Unlock()
	for id, cmd := range procs {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		delete(procs, id)
		running[id] = false
	}
}

func IsRunning(id string) bool {
	mu.Lock()
	defer mu.Unlock()
	return running[id]
}

func shell() string {
	if os.PathSeparator == '\\' {
		return "cmd"
	}
	return "sh"
}
func shellArg() string {
	if os.PathSeparator == '\\' {
		return "/C"
	}
	return "-c"
}
