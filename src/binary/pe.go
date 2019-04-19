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
	GetFilename() string
	GetSize() int64
	BuildXml(d *etree.Document) *etree.Element
}

type XmlBuildable interface {
	ToXml() *etree.Element
}

type XmlBuilderCallback func(item interface{}) *etree.Element

type Section struct {
	Name string
	Size uint64
	Flags string
}

type Dependency struct {
	Name string
}

type Flag struct {
	Name string
}

func (f *Flag) ToXml() *etree.Element {
	functionNode := etree.NewElement("Flag")
	functionNode.CreateText(f.Name)
	return functionNode
}

type Function struct {
	Name string
}

func (f *Function) ToXml() *etree.Element {
	functionNode := etree.NewElement("Function")
	functionNode.CreateText(f.Name)
	return functionNode
}

type PEBinary struct {
	Filename          string
	Architecture      string
	Signature         string
	Compiler          string
	Dependencies      []Dependency
	Size              int64
	Timestamp         int64
	Time              time.Time
	Flags             []Flag
	SectionNumber     uint16
	Sections          []*pe2.SectHeader
	ImportedFunctions []Function
	ExportedFunctions []Function
}

func (bin *PEBinary) GetFilename() string {
	return util.GetOptionalStringValue(bin.Filename, DEFAULT_VALUE)
}

func (bin *PEBinary) GetCompiler() string {
	return util.GetOptionalStringValue(bin.Compiler, DEFAULT_VALUE)
}

func (bin *PEBinary) GetSize() int64 {
	return bin.Size
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

func (bin *PEBinary) BuildXml(doc *etree.Document) *etree.Element {
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	root := doc.CreateElement("Binary")
	root.AddChild(util.BuildNodeWithText("Filename", bin.GetFilename()))
	root.AddChild(util.BuildNodeWithText("Architecture", bin.GetArchitecture()))
	root.AddChild(util.BuildNodeWithText("Signature", bin.GetSignature()))
	root.AddChild(util.BuildNodeWithText("Compiler", bin.GetCompiler()))

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
