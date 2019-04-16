package binary

import (
	"github.com/beevik/etree"
	"util"
)

type Binary interface {
	ToXml(d *etree.Document) *etree.Element
}

type Section struct {
	Idx int
	Name string
	Size int
}

type Dependency struct {
	Name string
}

type Flag struct {
	Name string
}

type PEBinary struct {
	Filename string
	Architecture string
	Dependencies []Dependency
	Flags []Flag
	Sections []Section
}

func (bin *PEBinary) GetArchitecture() string {
	return util.GetOptionalStringValue(bin.Architecture, "Any")
}

func (bin *PEBinary) ToXml(doc *etree.Document) *etree.Element {
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	root := doc.CreateElement("Binary")
	root.AddChild(util.BuildNodeWithText("Filename", bin.Filename))
	root.AddChild(util.BuildNodeWithText("Architecture", bin.GetArchitecture()))

	dependenciesElem := root.CreateElement("Dependencies")
	for _, dependency := range bin.Dependencies {
		dependencyElem := dependenciesElem.CreateElement("Dependency")
		dependencyElem.CreateText(dependency.Name)
	}

	flagsElem := root.CreateElement("Flags")
	for _, flag := range bin.Flags {
		flagElem := flagsElem.CreateElement("Flag")
		flagElem.CreateText(flag.Name)
	}

	return root
}
