package binary

import (
	"github.com/beevik/etree"
	"util"
)

const (
	DEFAULT_MANIFEST_VERSION = "1.0"
)

type JarBinary struct {
	PEBinary
	ManifestVersion string
	ClassPath []string
	BuildJdk string
	MainClass string
	BuiltBy string
	CreatedBy string
	JarAnalyzerTree *etree.Element
}

func (jar *JarBinary) GetManifestVersion() string {
	return util.GetOptionalStringValue(jar.ManifestVersion, DEFAULT_MANIFEST_VERSION)
}

func (jar *JarBinary) GetCreatedBy() string {
	return util.GetOptionalStringValue(jar.CreatedBy, DEFAULT_VALUE)
}

func (jar *JarBinary) GetBuildJdk() string {
	return util.GetOptionalStringValue(jar.BuildJdk, DEFAULT_VALUE)
}

func (jar *JarBinary) GetBuiltBy() string {
	return util.GetOptionalStringValue(jar.BuiltBy, DEFAULT_VALUE)
}

func (jar *JarBinary) GetMainClass() string {
	return util.GetOptionalStringValue(jar.MainClass, DEFAULT_VALUE)
}

func (jar *JarBinary) ToXml(doc *etree.Document) *etree.Element {
	root := jar.PEBinary.ToXml(doc)

	root.AddChild(util.BuildNodeWithText("ManifestVersion", jar.GetManifestVersion()))
	root.AddChild(util.BuildNodeWithText("CreatedBy", jar.GetCreatedBy()))
	root.AddChild(util.BuildNodeWithText("BuiltBy", jar.GetBuiltBy()))
	root.AddChild(util.BuildNodeWithText("BuildJdk", jar.GetBuildJdk()))
	root.AddChild(util.BuildNodeWithText("MainClass", jar.GetMainClass()))

	classPathsElement := root.CreateElement("ClassPaths")
	for _, cp := range jar.ClassPath {
		if cp == "" {
			continue
		}
		classPathsElement.AddChild(util.BuildNodeWithText("ClassPath", cp))
	}

	if len(classPathsElement.ChildElements()) == 0 {
		classPathsElement.CreateText(EMPTY)
	}

	for _, c := range jar.JarAnalyzerTree.ChildElements() {
		if len(c.ChildElements()) > 0 {
			root.AddChild(c)
		}
	}

	return root
}