package pe

import (
	"github.com/decomp/exp/bin"
	"github.com/decomp/exp/bin/pe"
	"github.com/micromagicman/binary-info/executable"
)

type DecompPE struct {
	PEFileInfo *bin.File
}

func (dp *DecompPE) GetName() string {
	return "decomp/pe"
}

func (dp *DecompPE) LoadFile(pathToExecutable string) bool {
	peFile, err := pe.ParseFile(pathToExecutable)
	if err != nil {
		return false
	}
	dp.PEFileInfo = peFile
	return true
}

func (dp *DecompPE) Process(e executable.Executable) {
	peFile := e.(*executable.PortableExecutable)
	peFile.Architecture = dp.PEFileInfo.Arch.String()
}