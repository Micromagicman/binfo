package wrapper

import (
	"binfo/executable"
	"binfo/os"
)

type ELFReaderUtil struct {
	BaseDump
}

func (eru *ELFReaderUtil) GetName() string {
	return "elfreader-util"
}

func (eru *ELFReaderUtil) LoadFile(pathToExecutable string) bool {
	command := os.Exec.ELFReaderCommand(pathToExecutable)
	stdOut, err := os.Exec.Execute(command)
	if err != nil {
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
	return Group(eru.Find("Operating System:\\s+([^.]+)"), 1)
}

func (eru *ELFReaderUtil) getFormat() string {
	return Group(eru.Find("Bit format:\\s+([^.]+)"), 1)
}

func (eru *ELFReaderUtil) getEndianess() string {
	return Group(eru.Find("Endianess:\\s+([^.]+)"), 1)
}

func (eru *ELFReaderUtil) getVersion() string {
	return Group(eru.Find("ELF Version:\\s+([^.]+)"), 1)
}

