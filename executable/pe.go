package executable

import (
	"binfo/util"
	"github.com/beevik/etree"
	pe2 "github.com/mewrev/pe"
)

type PortableExecutable struct {
	BaseExecutable
	ImExporter
	Addresses     map[string]string
	Libraries     []Library
	Flags         []Flag
	SectionNumber uint16
	Sections      []*pe2.SectHeader
}

func (bin *PortableExecutable) GetMagic() string {
	return "0x4D5A"
}

func (bin *PortableExecutable) BuildXml(doc *etree.Document) *etree.Element {
	root := BuildBaseBinaryInfo(bin, doc)
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

	if len(bin.Addresses) > 0 {
		addressesNode := root.CreateElement("Addresses")
		for name, address := range bin.Addresses {
			addressesNode.AddChild(util.BuildNodeWithText(name, address))
		}
	}

	bin.BuildImportsAndExports(root)

	if len(bin.Sections) > 0 {
		sectionsNode := root.CreateElement("Sections")
		for _, section := range bin.Sections {
			sectionNode := sectionsNode.CreateElement("Section")
			sectionNode.CreateElement("Name").CreateText(section.Name)
			sizeNode := sectionNode.CreateElement("Size")
			sizeNode.CreateAttr("unit", "bytes")
			sizeNode.CreateText(util.Uint32ToString(section.Size))
			sectionNode.CreateElement("Flags").CreateText(section.Flags.String())
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
