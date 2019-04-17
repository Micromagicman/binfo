package analyzer

import (
	"os"
	"os/exec"
)

const (
	ANALYZERS_PATH = "C:\\Users\\Admin\\Work\\binfo\\backend"
	DEFAULT_WORKING_DIR = "C:\\Users\\Admin"
	TEMPLATE_DIRECTORY = "C:\\Users\\Admin\\Work\\temp"
)

func Execute(command string) []byte {
	return ExecuteIn(command, DEFAULT_WORKING_DIR)
}

func ExecuteIn(command string, workingDir string) []byte {
	cmd := exec.Command("cmd", "/C", command)
	cmd.Dir = workingDir
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		return []byte{}
	}

	return stdoutStderr
}

func CreateTemplateDirectory() {
	if _, err := os.Stat(TEMPLATE_DIRECTORY); os.IsNotExist(err) {
		err := os.MkdirAll(TEMPLATE_DIRECTORY, os.ModePerm)
		if err != nil {
			panic("Error creating template directory")
		}
	}
}
