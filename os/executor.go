package os

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type CommandRunnable interface {
	GetWindowsCommand() string
	GetLinuxCommand() string
}

var Sep = string(os.PathSeparator)
var WorkingDir = "." + Sep
var BackendDir = WorkingDir + "binfo-backend"
var TemplateDir = WorkingDir + "temp"

func Execute(cr CommandRunnable, arguments... string) ([]byte, error) {
	cmd := exec.Command(cr.GetLinuxCommand(), arguments...)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cr.GetWindowsCommand() + " " + strings.Join(arguments, " "))
	}
	cmd.Dir = WorkingDir
	stdoutStderr, err := cmd.CombinedOutput()
	if nil != err {
		return nil, err
	}
	return stdoutStderr, nil
}