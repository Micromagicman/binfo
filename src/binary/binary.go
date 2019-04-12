package binary

import (
	"github.com/beevik/etree"
	"strconv"
)

type Binary interface {
	ToXml(d *etree.Document)
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

type JarBinary struct {
	PEBinary
	ManifestVersion string
	ClassPath []string
	BuildJdk string
	MainClass string
	BuildBy string
	CreatedBy string
}

func (bin *PEBinary) GetArchitecture() string {
	if bin.Architecture == "" {
		return "Any"
	}
	return bin.Architecture
}

func (bin *PEBinary) ToXml(doc *etree.Document) *etree.Element {
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	root := doc.CreateElement("Binary")
	filenameElem := root.CreateElement("Filename")
	filenameElem.CreateText(bin.Filename)

	bin.buildFlags(root)
	bin.buildDependencies(root)

	return root
}

func (jar *JarBinary) ToXml(doc *etree.Document) *etree.Document {
	root := jar.PEBinary.ToXml(doc)
	root.CreateElement("Jar")
	return nil
}

func (bin *PEBinary) buildArchitecture(root *etree.Element) {
	architectureElem := root.CreateElement("Architecture")
	architectureElem.CreateText(bin.GetArchitecture())
}

func (bin *PEBinary) buildDependencies(root *etree.Element) {
	dependenciesElem := root.CreateElement("Dependencies")
	for index, dependency := range bin.Dependencies {
		dependencyElem := dependenciesElem.CreateElement("Dependency")
		dependencyElem.CreateAttr("id", strconv.Itoa(index))
		dependencyElem.CreateText(dependency.Name)
	}
}

func (bin *PEBinary) buildFlags(root *etree.Element) {
	flagsElem := root.CreateElement("Flags")
	for index, flag := range bin.Flags {
		flagElem := flagsElem.CreateElement("Flag")
		flagElem.CreateAttr("id", strconv.Itoa(index))
		flagElem.CreateText(flag.Name)
	}
}
