package binary

import (
	"github.com/beevik/etree"
	"util"
)

const (
	DEFAULT_FORMAT    = "Executable"
	DEFAULT_ENDIANESS = "Little endian"
	DEFAULT_VERSION   = "Original ELF"
)

type ELFBinary struct {
	PEBinary
	Format          string
	Endianess       string
	Version         string
	SectionCount    uint16
	OperatingSystem string
	UnusedBytes     string
	Type            string
}

func (elf *ELFBinary) GetFormat() string {
	return util.GetOptionalStringValue(elf.Format, DEFAULT_FORMAT)
}

func (elf *ELFBinary) GetEndianess() string {
	return util.GetOptionalStringValue(elf.Endianess, DEFAULT_ENDIANESS)
}

func (elf *ELFBinary) GetVersion() string {
	return util.GetOptionalStringValue(elf.Version, DEFAULT_VERSION)
}

func (elf *ELFBinary) GetSectionCount() string {
	return util.GetOptionalStringValue("", DEFAULT_VALUE)
}

func (elf *ELFBinary) GetOperatingSystem() string {
	return util.GetOptionalStringValue(elf.OperatingSystem, DEFAULT_VALUE)
}

func (elf *ELFBinary) GetType() string {
	return util.GetOptionalStringValue(elf.Type, DEFAULT_VALUE)
}

func (elf *ELFBinary) ToXml(doc *etree.Document) *etree.Element {
	root := elf.PEBinary.ToXml(doc)
	root.AddChild(util.BuildNodeWithText("Format", elf.GetFormat()))
	root.AddChild(util.BuildNodeWithText("Endianess", elf.GetEndianess()))
	root.AddChild(util.BuildNodeWithText("ElfVersion", elf.GetVersion()))
	root.AddChild(util.BuildNodeWithText("OperatingSystem", elf.GetOperatingSystem()))
	root.AddChild(util.BuildNodeWithText("Type", elf.GetType()))
	return root
}
