package elf

import (
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/util"
)

type CDetect struct {
	Compiler string
}

func (cd *CDetect) GetWindowsCommand() string {
	return "call " + os.BackendDir + os.Sep + "cdetect.exe"
}

func (cd *CDetect) GetLinuxCommand() string {
	return os.BackendDir + os.Sep + "cdetect"
}

func (cd *CDetect) GetName() string {
	return "cdetect"
}

func (cd *CDetect) LoadFile(pathToExecutable string) bool {
	stdOut, err := os.Execute(cd, pathToExecutable)
	if nil != err {
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
