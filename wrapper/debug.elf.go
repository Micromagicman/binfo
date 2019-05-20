package wrapper

import (
	"binfo/executable"
	"debug/elf"
	"github.com/fatih/set"
)

type DebugElf struct {
	File *elf.File
}

func (de *DebugElf) GetName() string {
	return "debug/elf"
}

func (de *DebugElf) LoadFile(executablePath string) bool {
	file, err := elf.Open(executablePath)
	if err != nil {
		return false
	}
	de.File = file
	return true
}

func (de *DebugElf) Process(e executable.Executable) {
	elfFile, ok := e.(*executable.ExecutableLinkable)
	if !ok {
		return
	}
	elfFile.Libraries = set.New(set.NonThreadSafe)
	libraries, err := de.File.ImportedLibraries()
	if err == nil {
		for _, l := range libraries {
			elfFile.Libraries.Add(l)
		}
	}
	symbols, err := de.File.ImportedSymbols()
	if err == nil {
		imports := make([]executable.Function, len(symbols))
		for i, s := range symbols {
			elfFile.Libraries.Add(s.Library)
			imports[i] = executable.Function{s.Name, s.Library}
		}
		elfFile.Imports = imports
	}
}
