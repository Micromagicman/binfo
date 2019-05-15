package executable

import (
	"binfo/util"
	"github.com/beevik/etree"
	pe2 "github.com/mewrev/pe"
)

type PortableExecutable struct {
	BaseExecutable
	ImExporter
	Libraries     []Library
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

func (bin *PortableExecutable) GetOsVersion() string {
	return util.GetOptionalStringValue(bin.OsVersion, DEFAULT_VALUE)
}

func (bin *PortableExecutable) GetLinkerVersion() string {
	return util.GetOptionalStringValue(bin.LinkerVersion, DEFAULT_VALUE)
}

func (bin *PortableExecutable) GetMagic() string {
	return "0x4D5A"
}

func (bin *PortableExecutable) BuildXml(doc *etree.Document) *etree.Element {
	root := BuildBaseBinaryInfo(bin, doc)

	root.CreateElement("OsVersion").
		CreateText(bin.GetOsVersion())
	root.CreateElement("LinkerVersion").
		CreateText(bin.GetLinkerVersion())
	root.CreateElement("Checksum").
		CreateText(bin.Checksum)
	root.CreateElement("CodeRVA").
		CreateText(bin.CodeRVA)
	root.AddChild(buildSizeTag("CodeSize", bin.CodeSize))
	root.CreateElement("DataRVA").
		CreateText(bin.DataRVA)
	root.AddChild(buildSizeTag("DataSize", bin.DataSize))

	if len(bin.Libraries) > 0 {
		dependenciesNode := root.CreateElement("Libraries")
		for _, dependency := range bin.Libraries {
			dependencyNode := dependenciesNode.CreateElement("Library")
			dependencyNode.CreateText(dependency.Name)
		}
	}

	if len(bin.Flags) > 0 {
		flagsNode := root.CreateElement("Flags")
		for _, flag := range bin.Flags {
			flagNode := flagsNode.CreateElement("Flag")
			flagNode.CreateText(flag.Name)
		}
	}

	bin.BuildImportsAndExports(root)

	if len(bin.Sections) > 0 {
		sectionsNode := root.CreateElement("Sections")
		for _, section := range bin.Sections {
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
