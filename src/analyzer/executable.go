package analyzer

import (
	"strings"
)

func Objdump(binaryFilePath string, args ...string) *ObjDump {
	flagsString := "-" + strings.Join(args, "")
	stdOut := Execute("call " + ANALYZERS_PATH + "\\binutils\\objdump.exe " + binaryFilePath + " " + flagsString)
	return &ObjDump{string(stdOut)}
}

func PEDumper(binaryFilePath string) *ObjDump {
	stdOut := Execute("call " + ANALYZERS_PATH + "\\pedumper.exe " + binaryFilePath)
	return &ObjDump{string(stdOut)}
}