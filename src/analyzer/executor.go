package analyzer

import (
	"log"
	"os/exec"
	"runtime"
)

type Executor struct {
	DefaultWorkingDirectory string
	AnalyzersPath           string
	TemplateDirectory       string
	Sep                     string
	RunningCommands         map[string]string
}

func (e *Executor) Execute(command string) []byte {
	return e.ExecuteIn(command, e.DefaultWorkingDirectory)
}

func (e *Executor) ExecuteIn(command string, workingDir string) []byte {
	cmd := exec.Command(command)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	}

	cmd.Dir = workingDir
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
		return []byte{}
	}

	return stdoutStderr
}

func (e *Executor) ObjDumpCommand(binaryFilePath string, flags string) string {
	return e.RunningCommands["objdump"] + " " + binaryFilePath + " " + flags
}

func (e *Executor) PEDumperCommand(binaryFilePath string) string {
	return e.RunningCommands["pedumper"] + " " + binaryFilePath
}

func ExecutorFactory() *Executor {
	switch runtime.GOOS {
		case "windows": return createWindowsExecutor()
		default: return createLinuxExecutor()
	}
}

func createLinuxExecutor() *Executor {
	linExec := new(Executor)
	linExec.DefaultWorkingDirectory = "~"
	linExec.AnalyzersPath = "~/binfo/backend/linux"
	linExec.TemplateDirectory = "~/temp"
	linExec.Sep = "/"
	return linExec
}

func createWindowsExecutor() *Executor {
	winExec := new(Executor)
	winExec.DefaultWorkingDirectory = "C:\\Users\\Admin"
	winExec.AnalyzersPath = "C:\\Users\\Admin\\Work\\binfo\\backend\\windows"
	winExec.TemplateDirectory = "C:\\Users\\Admin\\Work\\temp"
	winExec.Sep = "\\"
	winExec.RunningCommands = map[string]string {
		"objdump": "call " + winExec.AnalyzersPath + "\\binutils\\objdump.exe",
		"pedumper": "call " + winExec.AnalyzersPath + "\\pedumper.exe",
		"jaranalyzer": "call " + winExec.AnalyzersPath + "\\jaranalyzer\\runxmlsummary.bat",
	}
	return winExec
}
