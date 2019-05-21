package elf

import (
	"binfo/executable"
	"github.com/decomp/exp/bin"
	"github.com/decomp/exp/bin/elf"
)

type DecompELF struct {
	ElfFile *bin.File
}

func (de *DecompELF) GetName() string {
	return "decomp/elf"
}

func (de *DecompELF) LoadFile(pathToExecutable string) bool {
	elfFile, err := elf.ParseFile(pathToExecutable)
	if err != nil {
		return false
	}
	de.ElfFile = elfFile
	return true
}

func (de *DecompELF) Process(e executable.Executable) {
	elfFile := e.(*executable.ExecutableLinkable)
	elfFile.Architecture = de.ElfFile.Arch.String()
	elfFile.Exports = de.ElfFile.Exports
}