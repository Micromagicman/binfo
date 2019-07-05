package analyzer

import (
	"github.com/micromagicman/binary-info/wrapper"
	"github.com/micromagicman/binary-info/wrapper/common"
	"github.com/micromagicman/binary-info/wrapper/elf"
	"github.com/micromagicman/binary-info/wrapper/jar"
	"github.com/micromagicman/binary-info/wrapper/pe"
)

type UtilitiesContainer struct {
	Common []wrapper.LibraryWrapper
	PE []wrapper.LibraryWrapper
	ELF []wrapper.LibraryWrapper
	JAR []wrapper.LibraryWrapper
}

func BuildUtilitiesContainer() *UtilitiesContainer {
	container := new(UtilitiesContainer)
	// Общие утилиты
	container.Common = []wrapper.LibraryWrapper {
		new(common.FileStat),
	}
	// Библиотеки для PE
	container.PE = []wrapper.LibraryWrapper{
		new(pe.DebugPe),
		new(pe.MemrevPE),
		new(pe.ObjDump),
		new(pe.DecompPE),
		new(pe.PEFile),
	}
	// Библиотеки для ELF
	container.ELF = []wrapper.LibraryWrapper{
		new(elf.DebugElf),
		new(elf.ELFReader),
		new(elf.ELFReaderUtil),
		new(elf.CDetect),
		new(elf.DecompELF),
	}
	// Библиотеки для Jar
	container.JAR = []wrapper.LibraryWrapper{
		new(jar.JarAnalyzer),
		new(jar.Tattletale),
	}
	return container
}
