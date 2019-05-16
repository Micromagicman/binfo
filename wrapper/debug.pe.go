package wrapper

import (
	"binfo/executable"
	"debug/pe"
	"github.com/fatih/set"
	"strings"
)

type DebugPe struct {
	File *pe.File
}

func CreateDebugPeWrapper(executablePath string) (*DebugPe, error) {
	file, err := pe.Open(executablePath)
	if err != nil {
		return nil, err
	}
	return &DebugPe{file}, nil
}

func (de *DebugPe) Process(pe *executable.PortableExecutable) {
	pe.Libraries = set.New(set.NonThreadSafe)
	libraries, err := de.File.ImportedLibraries()
	if err == nil {
		for _, l := range libraries {
			pe.Libraries.Add(l)
		}
	}
	symbols, err := de.File.ImportedSymbols()
	if err == nil {
		imports := make([]executable.Function, len(symbols))
		for i, s := range symbols {
			nameFrom := strings.Split(s, ":")
			pe.Libraries.Add(nameFrom[1])
			imports[i] = executable.Function{nameFrom[0], nameFrom[1]}
		}
		pe.Imports = imports
	}
}
