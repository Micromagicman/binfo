package binary

import (
	"binfo/util"
	"github.com/beevik/etree"
	"time"
)

type Binary interface {
	GetFilename() string
	GetArchitecture() string
	GetSize() int64
	GetTimestamp() int64
	GetDMY() string
	GetCompiler() string
	GetMagic() string
	GetProgrammingLanguage() string
	BuildXml(document *etree.Document) *etree.Element
}

const (
	EMPTY                = "Empty"
	DEFAULT_VALUE        = "Unknown"
	DEFAULT_ARCHITECTURE = "Any"
)

type XmlBuildable interface {
	ToXml() *etree.Element
}

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

type BaseBinary struct {
	Filename string
	Size int64
	Architecture string
	Timestamp int64
	Time time.Time
	Compiler string
	ProgrammingLanguage string
}

func (bin *BaseBinary) GetFilename() string {
	return util.GetOptionalStringValue(bin.Filename, DEFAULT_VALUE)
}

func (bin *BaseBinary) GetCompiler() string {
	return util.GetOptionalStringValue(bin.Compiler, DEFAULT_VALUE)
}

func (bin *BaseBinary) GetProgrammingLanguage() string {
	return util.GetOptionalStringValue(bin.ProgrammingLanguage, DEFAULT_VALUE)
}

func (bin *BaseBinary) GetSize() int64 {
	return bin.Size
}

func (bin *BaseBinary) GetArchitecture() string {
	return util.GetOptionalStringValue(bin.Architecture, DEFAULT_ARCHITECTURE)
}

func (bin *BaseBinary) GetDMY() string {
	return bin.Time.String()
}

func (bin *BaseBinary) GetTimestamp() int64 {
	return bin.Timestamp
}

func BuildBaseBinaryInfo(bin Binary, doc *etree.Document) *etree.Element {
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := doc.CreateElement("Binary")
	root.AddChild(util.BuildNodeWithText("Filename", bin.GetFilename()))
	root.AddChild(util.BuildNodeWithText("Magic", bin.GetMagic()))
	root.AddChild(util.BuildNodeWithText("Architecture", bin.GetArchitecture()))
	root.AddChild(util.BuildNodeWithText("Compiler", bin.GetCompiler()))
	root.AddChild(util.BuildNodeWithText("ProgrammingLanguage", bin.GetProgrammingLanguage()))

	if bin.GetSize() > 0 {
		sizeNode := root.CreateElement("Size")
		sizeNode.CreateAttr("unit", "bytes")
		sizeNode.CreateText(util.Int64ToString(bin.GetSize()))
	}

	if bin.GetTimestamp() > 0 {
		dateNode := root.CreateElement("CompilationDate")
		dateNode.CreateElement("Timestamp").CreateText(util.Int64ToString(bin.GetTimestamp()))
		dateNode.CreateElement("Date").CreateText(bin.GetDMY())
	}

	return root
}