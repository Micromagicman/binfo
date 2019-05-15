package wrapper

import (
	"binfo/executable"
	"github.com/yalue/elf_reader"
	"io/ioutil"
)

type ELFReader struct {
	File elf_reader.ELFFile
}

func CreateELFReaderWrapper(pathToElf string) (*ELFReader, error) {
	raw, e := ioutil.ReadFile(pathToElf)
	if e != nil {
		return nil, e
	}

	elf, e := elf_reader.ParseELFFile(raw)
	if e != nil {
		return nil, e
	}

	return &ELFReader{elf}, nil
}

func (er *ELFReader) Process(bin *executable.ExecutableLinkable) {
	bin.Sections = er.getSections()
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

//func (er *ELFReader) GetImportedFunctions() []executable.Function {
//	count := er.File.GetSectionCount()
//	functions := []executable.Function{}
//
//	for i := uint16(1); i < count; i++ {
//		_, functionNames, err := er.File.GetSymbols(i)
//		if err != nil {
//			continue
//		}
//
//		for _, functionName := range functionNames {
//			if functionName != "" && !strings.Contains(functionName, ".") && !strings.Contains(functionName, "@") {
//				functions = append(functions, executable.Function{functionName})
//			}
//		}
//	}
//
//	return functions
//}
