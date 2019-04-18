package analyzer

import (
	"strings"
)

func (a *Analyzer) ObjDump(binaryFilePath string, args ...string) *ObjDump {
	flagsString := "-" + strings.Join(args, "")
	command := a.Executor.ObjDumpCommand(binaryFilePath, flagsString)
	stdOut := a.Executor.Execute(command)
	return &ObjDump{string(stdOut)}
}

func (a *Analyzer) PEDumper(binaryFilePath string) *ObjDump {
	command := a.Executor.PEDumperCommand(binaryFilePath)
	stdOut := a.Executor.Execute(command)
	return &ObjDump{string(stdOut)}
}