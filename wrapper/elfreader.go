package wrapper

import (
	"binfo/binary"
	"fmt"
	"github.com/yalue/elf_reader"
	"io/ioutil"
	"strings"
)

type ELFReader struct {
	File elf_reader.ELFFile
}

func CreateELFReader(pathToElf string) (*ELFReader, error) {
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

func (er *ELFReader) GetSections() []binary.Section {
	count := er.File.GetSectionCount()
	sections := make([]binary.Section, int(count))

	for i := uint16(1); i < count; i++ {
		elfSectionHeader, err := er.File.GetSectionHeader(i)
		if err != nil {
			continue
		}

		section := binary.Section{}
		section.Size = uint64(elfSectionHeader.GetSize())
		sectionName, err := er.File.GetSectionName(i)
		if err != nil || sectionName == "" {
			fmt.Println(sectionName)
			sectionName = binary.DEFAULT_VALUE
		}

		section.Name = sectionName
		section.Flags = elfSectionHeader.GetFlags().String()
		sections[i] = section
	}

	return sections
}

func (er *ELFReader) GetImportedFunctions() []binary.Function {
	count := er.File.GetSectionCount()
	functions := []binary.Function{}

	for i := uint16(1); i < count; i++ {
		_, functionNames, err := er.File.GetSymbols(i)
		if err != nil {
			continue
		}

		for _, functionName := range functionNames {
			if functionName != "" && !strings.Contains(functionName, ".") && !strings.Contains(functionName, "@") {
				functions = append(functions, binary.Function{functionName})
			}
		}
	}

	return functions
}

