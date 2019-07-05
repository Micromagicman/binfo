package os

import (
	"fmt"
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
		fmt.Println(cr.GetWindowsCommand() + " " + strings.Join(arguments, " "))
		cmd = exec.Command("cmd", "/C", cr.GetWindowsCommand() + " " + strings.Join(arguments, " "))
	}
	cmd.Dir = WorkingDir
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return stdoutStderr, nil
}

func move(slice []int, to int) {
	for i, _ := range slice {
		temp := slice[i]
		slice[i] = slice[(i + to) % len(slice)]
		slice[(i + to) % len(slice)] = temp
	}
	fmt.Println(slice)
}