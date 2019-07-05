package executable

import (
	"github.com/beevik/etree"
	"github.com/fatih/set"
	pe2 "github.com/mewrev/pe"
	"github.com/micromagicman/binary-info/util"
)

type PortableExecutable struct {
	BaseExecutable
	ImExporter
	Libraries     set.Interface
	Flags         []Flag
	SectionNumber uint16
	Sections      []*pe2.SectHeader
	LinkerVersion string
	OsVersion     string
	Checksum      string
	CodeRVA       string
	CodeSize      string
	DataRVA       string
	DataSize      string
}

func (pe *PortableExecutable) GetType() string {
	return "Windows Portable Executable"
}

func (pe *PortableExecutable) GetOsVersion() string {
	return util.GetOptionalStringValue(pe.OsVersion, DEFAULT_VALUE)
}

func (pe *PortableExecutable) GetLinkerVersion() string {
	return util.GetOptionalStringValue(pe.LinkerVersion, DEFAULT_VALUE)
}

func (pe *PortableExecutable) GetMagic() string {
	return "0x4D5A"
}

func (pe *PortableExecutable) BuildXml(doc *etree.Document) *etree.Element {
	root := BuildBaseBinaryInfo(pe, doc)

	root.CreateElement("OsVersion").
		CreateText(pe.GetOsVersion())
	root.CreateElement("LinkerVersion").
		CreateText(pe.GetLinkerVersion())
	root.CreateElement("Checksum").
		CreateText(pe.Checksum)
	root.CreateElement("CodeRVA").
		CreateText(pe.CodeRVA)
	root.AddChild(buildSizeTag("CodeSize", pe.CodeSize))
	root.CreateElement("DataRVA").
		CreateText(pe.DataRVA)
	root.AddChild(buildSizeTag("DataSize", pe.DataSize))

	if pe.Libraries.Size() > 0 {
		root.AddChild(buildLibraries(pe.Libraries))
	}

	if len(pe.Flags) > 0 {
		flagsNode := root.CreateElement("Flags")
		for _, flag := range pe.Flags {
			flagNode := flagsNode.CreateElement("Flag")
			flagNode.CreateText(flag.Name)
		}
	}

	pe.BuildImportsAndExports(root)

	if len(pe.Sections) > 0 {
		sectionsNode := root.CreateElement("Sections")
		for _, section := range pe.Sections {
			sectionNode := sectionsNode.CreateElement("Section")
			sectionNode.CreateElement("Name").
				CreateText(util.GetOptionalStringValue(section.Name, DEFAULT_VALUE))
			sectionNode.AddChild(buildSizeTag("Size", util.Uint32ToString(section.Size)))
			sectionNode.CreateElement("Flags").
				CreateText(section.Flags.String())
		}
	}

	return root
}

func arrayOf(name string, items []XmlBuildable) *etree.Element {
	wrapperNode := etree.NewElement(name)
	for _, item := range items {
		wrapperNode.AddChild(item.ToXml())
	}
	return wrapperNode
}
