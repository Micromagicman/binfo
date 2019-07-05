package elf

import (
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/wrapper"
)

type ELFReaderUtil struct {
	wrapper.BaseDump
}

func (eru *ELFReaderUtil) GetWindowsCommand() string {
	return "call " + os.BackendDir + os.Sep + "elfreader.exe"
}

func (eru *ELFReaderUtil) GetLinuxCommand() string {
	return os.BackendDir + os.Sep + "elfreader"
}

func (eru *ELFReaderUtil) GetName() string {
	return "elfreader-util"
}

func (eru *ELFReaderUtil) LoadFile(pathToExecutable string) bool {
	stdOut, err := os.Execute(eru, pathToExecutable)
	if nil != err {
		return false
	}
	eru.Content = string(stdOut)
	return true
}

func (eru *ELFReaderUtil) Process(e executable.Executable) {
	elfFile := e.(*executable.ExecutableLinkable)
	elfFile.Format = eru.getFormat()
	elfFile.Version = eru.getVersion()
	elfFile.Endianess = eru.getEndianess()
	elfFile.OperatingSystem = eru.getOperatingSystem()
}

func (eru *ELFReaderUtil) getOperatingSystem() string {
	return wrapper.Group(eru.Find("Operating System:\\s+([^.]+)"), 1)
}

func (eru *ELFReaderUtil) getFormat() string {
	return wrapper.Group(eru.Find("Bit format:\\s+([^.]+)"), 1)
}

func (eru *ELFReaderUtil) getEndianess() string {
	return wrapper.Group(eru.Find("Endianess:\\s+([^.]+)"), 1)
}

func (eru *ELFReaderUtil) getVersion() string {
	return wrapper.Group(eru.Find("ELF Version:\\s+([^.]+)"), 1)
}

