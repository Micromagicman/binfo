package binary

import (
	"github.com/beevik/etree"
	pe2 "github.com/mewrev/pe"
	"time"
	"util"
)

const (
	EMPTY                = "Empty"
	DEFAULT_VALUE        = "Unknown"
	DEFAULT_ARCHITECTURE = "Any"
)

type Binary interface {
	ToXml(d *etree.Document) *etree.Element
}

type XmlBuilderCallback func(item interface{}) *etree.Element

type Section struct {
	Idx  int
	Name string
	Size int
}

type Dependency struct {
	Name string
}

type Flag struct {
	Name string
}

type Function struct {
	Name string
}

type PEBinary struct {
	Filename          string
	Architecture      string
	Signature         string
	Dependencies      []Dependency
	Size              int64
	Timestamp         int64
	Time              time.Time
	Flags             []Flag
	SectionNumber     uint16
	Sections          []*pe2.SectHeader
	ImportedFunctions []Function
}

func (bin *PEBinary) GetFilename() string {
	return util.GetOptionalStringValue(bin.Filename, DEFAULT_VALUE)
}

func (bin *PEBinary) GetArchitecture() string {
	return util.GetOptionalStringValue(bin.Architecture, DEFAULT_ARCHITECTURE)
}

func (bin *PEBinary) GetSignature() string {
	return util.GetOptionalStringValue(bin.Signature, DEFAULT_VALUE)
}

func (bin *PEBinary) GetDMY() string {
	return bin.Time.String()
}

func (bin *PEBinary) ToXml(doc *etree.Document) *etree.Element {
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	root := doc.CreateElement("Binary")
	root.AddChild(util.BuildNodeWithText("Filename", bin.GetFilename()))
	root.AddChild(util.BuildNodeWithText("Architecture", bin.GetArchitecture()))
	root.AddChild(util.BuildNodeWithText("Signature", bin.GetSignature()))

	if bin.Size > 0 {
		sizeNode := root.CreateElement("Size")
		sizeNode.CreateAttr("unit", "bytes")
		sizeNode.CreateText(util.Int64ToString(bin.Size))
	}

	if bin.Timestamp > 0 {
		dateNode := root.CreateElement("CompilationDate")
		dateNode.CreateElement("Timestamp").CreateText(util.Int64ToString(bin.Timestamp))
		dateNode.CreateElement("Date").CreateText(bin.GetDMY())
	}

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

	if len(bin.ImportedFunctions) > 0 {
		importedFunctionsNode := root.CreateElement("ImportedFunctions")
		for _, function := range bin.ImportedFunctions {
			funcNode := importedFunctionsNode.CreateElement("Function")
			funcNode.CreateText(function.Name)
		}
	}

	return root
}

//func arrayOf(name string, items []interface{}, callback XmlBuilderCallback) *etree.Element {
//	if len(items) > 0 {
//		return nil
//	}
//
//	parent := etree.NewElement(name)
//	for _, item := range items {
//		parent.AddChild(callback(item))
//	}
//
//	return parent
//}
