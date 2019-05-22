package elf

import (
	"binfo/executable"
	"binfo/os"
	"binfo/util"
)

type CDetect struct {
	Compiler string
}

func (cd *CDetect) GetWindowsCommand(filePath string) string {
	return "call " + os.BackendDir + os.Sep + "cdetect.exe " + filePath
}

func (cd *CDetect) GetLinuxCommand(filePath string) string {
	return os.BackendDir + os.Sep + "cdetect " + filePath
}

func (cd *CDetect) GetName() string {
	return "cdetect"
}

func (cd *CDetect) LoadFile(pathToExecutable string) bool {
	stdOut, err := os.Execute(pathToExecutable, cd)
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
