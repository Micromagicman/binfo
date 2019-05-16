package executable

import (
	"binfo/util"
	"github.com/decomp/exp/bin"
	"time"

	"github.com/beevik/etree"
)

type Executable interface {
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

type ImExporter struct {
	Imports []Function
	Exports map[bin.Address]string
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
	Name  string
	Size  uint64
	Flags string
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
	From string
}

func (f *Function) ToXml() *etree.Element {
	functionNode := etree.NewElement("Function")
	functionNode.CreateAttr("from", util.GetOptionalStringValue(f.From, "?"))
	functionNode.CreateText(f.Name)
	return functionNode
}

type BaseExecutable struct {
	Filename            string
	Size                int64
	Architecture        string
	Timestamp           int64
	Time                time.Time
	Compiler            string
	ProgrammingLanguage string
}

func (bin *BaseExecutable) GetFilename() string {
	return util.GetOptionalStringValue(bin.Filename, DEFAULT_VALUE)
}

func (bin *BaseExecutable) GetCompiler() string {
	return util.GetOptionalStringValue(bin.Compiler, DEFAULT_VALUE)
}

func (bin *BaseExecutable) GetProgrammingLanguage() string {
	return util.GetOptionalStringValue(bin.ProgrammingLanguage, DEFAULT_VALUE)
}

func (bin *BaseExecutable) GetSize() int64 {
	return bin.Size
}

func (bin *BaseExecutable) GetArchitecture() string {
	return util.GetOptionalStringValue(bin.Architecture, DEFAULT_ARCHITECTURE)
}

func (bin *BaseExecutable) GetDMY() string {
	return bin.Time.String()
}

func (bin *BaseExecutable) GetTimestamp() int64 {
	return bin.Timestamp
}

func (bin *BaseExecutable) GetMagic() string {
	return "Unknown"
}

func (bin *ImExporter) BuildImportsAndExports(root *etree.Element) {
	if len(bin.Imports) > 0 {
		listNode := root.CreateElement("Imports")
		for _, fn := range bin.Imports {
			listNode.AddChild(fn.ToXml())
		}
	}
	if len(bin.Exports) > 0 {
		listNode := root.CreateElement("Exports")
		for address, name := range bin.Exports {
			funcNode := listNode.CreateElement("Function")
			funcNode.CreateElement("Address").CreateText(address.String())
			funcNode.CreateElement("Name").CreateText(name)
		}
	}
}

func BuildBaseBinaryInfo(bin Executable, doc *etree.Document) *etree.Element {
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := doc.CreateElement("Executable")
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

	dateNode := root.CreateElement("CompileTime")
	unixTimeNode := dateNode.CreateElement("Unix")
	dateTimeNode := dateNode.CreateElement("DateTime")
	if bin.GetTimestamp() > 0 {
		unixTimeNode.CreateText(util.Int64ToString(bin.GetTimestamp()))
		dateTimeNode.CreateText(bin.GetDMY())
	} else {
		unixTimeNode.CreateText(DEFAULT_VALUE)
		dateTimeNode.CreateText(DEFAULT_VALUE)
	}

	return root
}

func buildSizeTag(name string, value string) *etree.Element {
	sizeTag := etree.NewElement(name)
	sizeTag.CreateAttr("unit", "bytes")
	sizeTag.CreateText(value)
	return sizeTag
}
