package wrapper

import (
	"binfo/executable"
	"binfo/os"
	"binfo/util"
)

type CDetect struct {
	Compiler string
}

func (cd *CDetect) GetName() string {
	return "cdetect"
}

func (cd *CDetect) LoadFile(pathToExecutable string) bool {
	command := os.Exec.CDetectCommand(pathToExecutable)
	stdOut, err := os.Exec.Execute(command)
	if err != nil {
		return false
	}
	cd.Compiler = string(stdOut)
	return true
}

func (cd *CDetect) Process(e executable.Executable) {
	elfFile := e.(*executable.ExecutableLinkable)
	elfFile.Compiler = cd.Compiler
	elfFile.ProgrammingLanguage = util.GetLanguageByCompiler(cd.Compiler)
}
