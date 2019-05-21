package pe

import (
	"binfo/executable"
	"debug/pe"
	"github.com/fatih/set"
	"strings"
)

type DebugPe struct {
	File *pe.File
}

func (dpe *DebugPe) GetName() string {
	return "debug/pe"
}

func (dpe *DebugPe) LoadFile(pathToExecutable string) bool {
	file, err := pe.Open(pathToExecutable)
	if err != nil {
		return false
	}
	dpe.File = file
	return true
}

func (dpe *DebugPe) Process(e executable.Executable) {
	peFile := e.(*executable.PortableExecutable)
	peFile.Libraries = set.New(set.NonThreadSafe)
	libraries, err := dpe.File.ImportedLibraries()
	if err == nil {
		for _, l := range libraries {
			peFile.Libraries.Add(l)
		}
	}
	symbols, err := dpe.File.ImportedSymbols()
	if err == nil {
		imports := make([]executable.Function, len(symbols))
		for i, s := range symbols {
			nameFrom := strings.Split(s, ":")
			peFile.Libraries.Add(nameFrom[1])
			imports[i] = executable.Function{nameFrom[0], nameFrom[1]}
		}
		peFile.Imports = imports
	}
}
