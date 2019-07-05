package elf

import (
	"github.com/decomp/exp/bin"
	"github.com/decomp/exp/bin/elf"
	"github.com/micromagicman/binary-info/executable"
)

type DecompELF struct {
	ElfFile *bin.File
}

func (de *DecompELF) GetName() string {
	return "decomp/elf"
}

func (de *DecompELF) LoadFile(pathToExecutable string) bool {
	elfFile, err := elf.ParseFile(pathToExecutable)
	if nil != err {
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