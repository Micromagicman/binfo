package analyzer

import (
	"binfo/wrapper"
	"binfo/wrapper/common"
	"binfo/wrapper/elf"
	"binfo/wrapper/jar"
	"binfo/wrapper/pe"
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
	// Библиотеки для ProcessJar
	container.JAR = []wrapper.LibraryWrapper{
		new(jar.Tattletale),
		new(jar.JarAnalyzer),
	}
	return container
}
