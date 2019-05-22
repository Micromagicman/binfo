package os

import (
	"os"
	"os/exec"
	"runtime"
)

type CommandRunnable interface {
	GetWindowsCommand(filePath string) string
	GetLinuxCommand(filePath string) string
}

var Sep = string(os.PathSeparator)
var WorkingDir = "." + Sep
var BackendDir = WorkingDir + "backend"
var TemplateDir = WorkingDir + "temp"

func Execute(filePath string, cr CommandRunnable) ([]byte, error) {
	cmd := exec.Command(cr.GetLinuxCommand(filePath))
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cr.GetWindowsCommand(filePath))
	}
	cmd.Dir = WorkingDir
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return stdoutStderr, nil
}