package binary

import (
	"binfo/util"
	"github.com/beevik/etree"
	"github.com/mewrev/pe"
)

const (
	DEFAULT_MANIFEST_VERSION = "1.0"
	DEFAULT_JAVA_COMPILER = "Javac"
)

type JarBinary struct {
	BaseBinary
	Dependencies    []Dependency
	Flags           []Flag
	Sections        []*pe.SectHeader
	ManifestVersion string
	ClassPath       []string
	BuildJdk        string
	MainClass       string
	BuiltBy         string
	CreatedBy       string
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

func (jar *JarBinary) GetCompiler() string {
	return util.GetOptionalStringValue(jar.Compiler, DEFAULT_JAVA_COMPILER)
}

func (jar *JarBinary) GetMagic() string {
	return "0x504B0304"
}

func (jar *JarBinary) BuildXml(doc *etree.Document) *etree.Element {
	root := BuildBaseBinaryInfo(jar, doc)
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

	if jar.JarAnalyzerTree != nil {
		for _, c := range jar.JarAnalyzerTree.ChildElements() {
			if len(c.ChildElements()) > 0 {
				root.AddChild(c)
			}
		}
	}

	return root
}
