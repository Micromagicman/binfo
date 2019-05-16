package analyzer

import (
	"os/exec"
	"path/filepath"
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
		return []byte{}, nil
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

func (e *Executor) TattletaleCommand(binaryFilePath string) string {
	return e.RunningCommands["tattletale"] + " " + filepath.Dir(binaryFilePath) + " " + e.TemplateDirectory
}

func (e *Executor) CDetectCommand(binaryFilePath string) string {
	return e.RunningCommands["cdetect"] + " " + binaryFilePath
}

func ExecutorFactory() *Executor {
	switch runtime.GOOS {
	case "windows":
		return createWindowsExecutor()
	default:
		return createLinuxExecutor()
	}
}

func createLinuxExecutor() *Executor {
	linExec := new(Executor)
	linExec.DefaultWorkingDirectory = "./"
	linExec.AnalyzersPath = linExec.DefaultWorkingDirectory + "backend/linux/"
	linExec.TemplateDirectory = linExec.DefaultWorkingDirectory + "temp/"
	linExec.Sep = "/"
	linExec.RunningCommands = map[string]string{
		"objdump":     linExec.AnalyzersPath + "objdump",
		"jaranalyzer": linExec.AnalyzersPath + "jaranalyzer/runxmlsummary",
		"elfreader":   linExec.AnalyzersPath + "elfreader",
		"cdetect":     linExec.AnalyzersPath + "cdetect",
		"tattletale":  "java -jar " + linExec.AnalyzersPath + "tattletale1.2.0.jar",
	}
	return linExec
}

func createWindowsExecutor() *Executor {
	winExec := new(Executor)
	winExec.DefaultWorkingDirectory = ".\\"
	winExec.AnalyzersPath = winExec.DefaultWorkingDirectory + "backend\\windows\\"
	winExec.TemplateDirectory = winExec.DefaultWorkingDirectory + "temp\\"
	winExec.Sep = "\\"
	winExec.RunningCommands = map[string]string{
		"objdump":     "call " + winExec.AnalyzersPath + "objdump.exe",
		"jaranalyzer": "call " + winExec.AnalyzersPath + "jaranalyzer\\runxmlsummary.bat",
		"elfreader":   "call " + winExec.AnalyzersPath + "elfreader.exe",
		"cdetect":     "call " + winExec.AnalyzersPath + "cdetect.exe",
		"tattletale":  "java -jar " + winExec.AnalyzersPath + "tattletale1.2.0.jar",
	}
	return winExec
}
