package wrapper

import (
	"binfo/executable"
	"github.com/yalue/elf_reader"
	"io/ioutil"
)

type ELFReader struct {
	File elf_reader.ELFFile
}

func (er *ELFReader) GetName() string {
	return "ElfReader"
}

func (er *ELFReader) LoadFile(pathToExecutable string) bool {
	raw, err := ioutil.ReadFile(pathToExecutable)
	if err != nil {
		return false
	}
	elf, err := elf_reader.ParseELFFile(raw)
	if err != nil {
		return false
	}
	er.File = elf
	return true
}

func (er *ELFReader) Process(e executable.Executable) {
	elfFile := e.(*executable.ExecutableLinkable)
	elfFile.Sections = er.getSections()
}

func (er *ELFReader) getSections() []executable.Section {
	count := er.File.GetSectionCount()
	sections := make([]executable.Section, int(count))

	for i := uint16(1); i < count; i++ {
		elfSectionHeader, err := er.File.GetSectionHeader(i)
		if err != nil {
			continue
		}

		section := executable.Section{}
		section.Size = uint64(elfSectionHeader.GetSize())
		sectionName, err := er.File.GetSectionName(i)
		if err != nil || sectionName == "" {
			sectionName = executable.DEFAULT_VALUE
		}

		section.Name = sectionName
		section.Flags = elfSectionHeader.GetFlags().String()
		sections[i] = section
	}

	return sections
}