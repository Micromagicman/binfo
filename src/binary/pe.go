package binary

import (
	"github.com/beevik/etree"
	pe2 "github.com/mewrev/pe"
	"util"
)

type PEBinary struct {
	BaseBinary
	Architecture      string
	Addresses         map[string]string
	Dependencies      []Dependency
	Flags             []Flag
	SectionNumber     uint16
	Sections          []*pe2.SectHeader
	ImportedFunctions []Function
	ExportedFunctions []Function
}

func (bin *PEBinary) GetMagic() string {
	return "0x4D5A"
}

func (bin *PEBinary) BuildXml(doc *etree.Document) *etree.Element {
	root := BuildBaseBinaryInfo(bin, doc)
	if len(bin.Dependencies) > 0 {
		dependenciesNode := root.CreateElement("Dependencies")
		for _, dependency := range bin.Dependencies {
			dependencyNode := dependenciesNode.CreateElement("Dependency")
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

	//root.AddChild(arrayOf("ImportedFunctions", bin.ImportedFunctions))
	//root.AddChild(arrayOf("ExportedFunction", bin.ExportedFunctions))
	if len(bin.ImportedFunctions) > 0 {
		importedFunctionsNode := root.CreateElement("ImportedFunctions")
		for _, function := range bin.ImportedFunctions {
			funcNode := importedFunctionsNode.CreateElement("Function")
			funcNode.CreateText(function.Name)
		}
	}

	if len(bin.ExportedFunctions) > 0 {
		importedFunctionsNode := root.CreateElement("ExportedFunctions")
		for _, function := range bin.ExportedFunctions {
			funcNode := importedFunctionsNode.CreateElement("Function")
			funcNode.CreateText(function.Name)
		}
	}

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
