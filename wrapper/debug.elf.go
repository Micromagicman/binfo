package wrapper

import (
	"binfo/executable"
	"debug/elf"
	"github.com/fatih/set"
)

type DebugElf struct {
	File *elf.File
}

func CreateDebugElfWrapper(executablePath string) (*DebugElf, error) {
	file, err := elf.Open(executablePath)
	if err != nil {
		return nil, err
	}
	return &DebugElf{file}, nil
}

func (de *DebugElf) Process(elf *executable.ExecutableLinkable) {
	elf.Libraries = set.New(set.NonThreadSafe)
	libraries, err := de.File.ImportedLibraries()
	if err == nil {
		for _, l := range libraries {
			elf.Libraries.Add(l)
		}
	}
	symbols, err := de.File.ImportedSymbols()
	if err == nil {
		imports := make([]executable.Function, len(symbols))
		for i, s := range symbols {
			elf.Libraries.Add(s.Library)
			imports[i] = executable.Function{s.Name, s.Library}
		}
		elf.Imports = imports
	}
}
