package executable

import (
	"github.com/beevik/etree"
	"github.com/fatih/set"
	"github.com/micromagicman/binary-info/util"
)

const (
	DEFAULT_FORMAT    = "Executable"
	DEFAULT_ENDIANESS = "Little endian"
	DEFAULT_VERSION   = "Original ELF"
)

type ExecutableLinkable struct {
	BaseExecutable
	ImExporter
	Libraries       set.Interface
	Format          string
	Endianess       string
	Version         string
	SectionCount    uint16
	OperatingSystem string
	UnusedBytes     string
	Sections        []Section
}

func (elf *ExecutableLinkable) GetType() string {
	return "Executable Linkable"
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

	if elf.Libraries.Size() > 0 {
		root.AddChild(buildLibraries(elf.Libraries))
	}

	elf.BuildImportsAndExports(root)

	if len(elf.Sections) > 0 {
		sectionsNode := root.CreateElement("Sections")
		for _, s := range elf.Sections {
			sectionNode := sectionsNode.CreateElement("Section")
			sectionNode.CreateElement("Name").CreateText(s.Name)
			sectionNode.AddChild(buildSizeTag("Size", util.UInt64ToString(s.Size)))
			sectionNode.CreateElement("Flags").CreateText(s.Flags)
		}
	}

	return root
}
