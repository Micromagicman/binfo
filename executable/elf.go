package executable

import (
	"binfo/util"

	"github.com/beevik/etree"
)

const (
	DEFAULT_FORMAT    = "Executable"
	DEFAULT_ENDIANESS = "Little endian"
	DEFAULT_VERSION   = "Original ELF"
)

type ExecutableLinkable struct {
	BaseExecutable
	Format            string
	Endianess         string
	Version           string
	SectionCount      uint16
	OperatingSystem   string
	UnusedBytes       string
	Type              string
	ImportedFunctions []Function
	Sections          []Section
}

func (elf *ExecutableLinkable) GetFormat() string {
	return util.GetOptionalStringValue(elf.Format, DEFAULT_FORMAT)
}

func (elf *ExecutableLinkable) GetEndianess() string {
	return util.GetOptionalStringValue(elf.Endianess, DEFAULT_ENDIANESS)
}

func (elf *ExecutableLinkable) GetVersion() string {
	return util.GetOptionalStringValue(elf.Version, DEFAULT_VERSION)
}

func (elf *ExecutableLinkable) GetSectionCount() string {
	return util.GetOptionalStringValue("", DEFAULT_VALUE)
}

func (elf *ExecutableLinkable) GetOperatingSystem() string {
	return util.GetOptionalStringValue(elf.OperatingSystem, DEFAULT_VALUE)
}

func (elf *ExecutableLinkable) GetType() string {
	return util.GetOptionalStringValue(elf.Type, DEFAULT_VALUE)
}

func (elf *ExecutableLinkable) GetMagic() string {
	return "0x7F454C46 (ELF)"
}

func (elf *ExecutableLinkable) BuildXml(doc *etree.Document) *etree.Element {
	root := BuildBaseBinaryInfo(elf, doc)
	root.AddChild(util.BuildNodeWithText("Format", elf.GetFormat()))
	root.AddChild(util.BuildNodeWithText("Endianess", elf.GetEndianess()))
	root.AddChild(util.BuildNodeWithText("ElfVersion", elf.GetVersion()))
	root.AddChild(util.BuildNodeWithText("OperatingSystem", elf.GetOperatingSystem()))
	root.AddChild(util.BuildNodeWithText("Type", elf.GetType()))

	if len(elf.ImportedFunctions) > 0 {
		importedFunctionsNode := root.CreateElement("ImportedFunctions")
		for _, function := range elf.ImportedFunctions {
			funcNode := importedFunctionsNode.CreateElement("Function")
			funcNode.CreateText(function.Name)
		}
	}

	if len(elf.Sections) > 0 {
		sectionsNode := root.CreateElement("Sections")
		for _, s := range elf.Sections {
			sectionNode := sectionsNode.CreateElement("Section")
			sectionNode.CreateElement("Name").CreateText(s.Name)
			sizeNode := sectionNode.CreateElement("Size")
			sizeNode.CreateAttr("unit", "bytes")
			sizeNode.CreateText(util.UInt64ToString(s.Size))
			sectionNode.CreateElement("Flags").CreateText(s.Flags)
		}
	}

	return root
}
