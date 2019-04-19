package analyzer

import (
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

func (e *Executor) Execute(command string) ([]byte, error) {
	return e.ExecuteIn(command, e.DefaultWorkingDirectory)
}

func (e *Executor) ExecuteIn(command string, workingDir string) ([]byte, error) {
	cmd := exec.Command(command)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	}

	cmd.Dir = workingDir
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	return stdoutStderr, nil
}

func (e *Executor) ObjDumpCommand(binaryFilePath string, flags string) string {
	return e.RunningCommands["objdump"] + " " + binaryFilePath + " " + flags
}

func (e *Executor) PEDumperCommand(binaryFilePath string) string {
	return e.RunningCommands["pedumper"] + " " + binaryFilePath
}

func (e *Executor) ELFReaderCommand(binaryFilePath string) string {
	return e.RunningCommands["elfreader"] + " " + binaryFilePath
}

func (e *Executor) ELFInfoCommand(binaryFilePath string) string {
	return e.RunningCommands["elfinfo"] + " -c " + binaryFilePath
}

func ExecutorFactory() *Executor {
	switch runtime.GOOS {
		case "windows": return createWindowsExecutor()
		default: return createLinuxExecutor()
	}
}

func createLinuxExecutor() *Executor {
	linExec := new(Executor)
	linExec.DefaultWorkingDirectory = "./"
	linExec.AnalyzersPath = linExec.DefaultWorkingDirectory + "backend/linux/"
	linExec.TemplateDirectory = linExec.DefaultWorkingDirectory + "temp/"
	linExec.Sep = "/"
	linExec.RunningCommands = map[string]string {
		"objdump": "objdump",
		"pedumper": linExec.AnalyzersPath + "pedumper",
		"jaranalyzer": linExec.AnalyzersPath + "jaranalyzer/runxmlsummary",
		"elfreader": linExec.AnalyzersPath + "elfreader",
		"readelf": "readelf",
	}
	return linExec
}

func createWindowsExecutor() *Executor {
	winExec := new(Executor)
	winExec.DefaultWorkingDirectory = ".\\"
	winExec.AnalyzersPath = winExec.DefaultWorkingDirectory + "backend\\windows\\"
	winExec.TemplateDirectory = winExec.DefaultWorkingDirectory +  "temp\\"
	winExec.Sep = "\\"
	winExec.RunningCommands = map[string]string {
		"objdump": "call " + winExec.AnalyzersPath + "objdump.exe",
		"pedumper": "call " + winExec.AnalyzersPath + "pedumper.exe",
		"jaranalyzer": "call " + winExec.AnalyzersPath + "jaranalyzer\\runxmlsummary.bat",
		"elfreader": "call " + winExec.AnalyzersPath + "elfreader.exe",
		"readelf": "call " + winExec.AnalyzersPath + "readelf.exe",
		"elfinfo": "call " + winExec.AnalyzersPath + "elfinfo.exe",
	}
	return winExec
}
