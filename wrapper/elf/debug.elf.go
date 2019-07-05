package elf

import (
	"debug/elf"
	"github.com/fatih/set"
	"github.com/micromagicman/binary-info/executable"
)

type DebugElf struct {
	File *elf.File
}

func (de *DebugElf) GetName() string {
	return "debug/elf"
}

func (de *DebugElf) LoadFile(executablePath string) bool {
	file, err := elf.Open(executablePath)
	if nil != err {
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
	if nil == err {
		for _, l := range libraries {
			elfFile.Libraries.Add(l)
		}
	}
	symbols, err := de.File.ImportedSymbols()
	if nil == err {
		imports := make([]executable.Function, len(symbols))
		for i, s := range symbols {
			elfFile.Libraries.Add(s.Library)
			imports[i] = executable.Function{s.Name, s.Library}
		}
		elfFile.Imports = imports
	}
}
