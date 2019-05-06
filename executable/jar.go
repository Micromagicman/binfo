package executable

import (
	"binfo/util"
	"strings"

	"github.com/beevik/etree"
)

const (
	DEFAULT_JAVA_COMPILER = "Javac"
)

type JarManifest map[string]string

type JarExecutable struct {
	BaseExecutable
	Requires        []string
	Provides        []string
	Manifest        JarManifest
	JarAnalyzerTree *etree.Element
}

func (jar *JarExecutable) GetCompiler() string {
	return util.GetOptionalStringValue(jar.Compiler, DEFAULT_JAVA_COMPILER)
}

func (jar *JarExecutable) GetMagic() string {
	return "0x504B0304"
}

func (jar *JarExecutable) CreateManifest() *etree.Element {
	manifestNode := etree.NewElement("Manifest")
	for key, value := range jar.Manifest {
		xmlKey := strings.ReplaceAll(key, "-", "")
		manifestNode.AddChild(util.BuildNodeWithText(xmlKey, value))
	}
	return manifestNode
}

func (jar *JarExecutable) BuildXml(doc *etree.Document) *etree.Element {
	root := BuildBaseBinaryInfo(jar, doc)
	root.AddChild(jar.CreateManifest())

	if jar.JarAnalyzerTree != nil {
		for _, c := range jar.JarAnalyzerTree.ChildElements() {
			if len(c.ChildElements()) > 0 {
				root.AddChild(c)
			}
		}
	}

	if len(jar.Requires) > 0 {
		requiresNode := root.CreateElement("Requires")
		for _, r := range jar.Requires {
			requiresNode.AddChild(util.BuildNodeWithText("Class", r))
		}
	}

	if len(jar.Provides) > 0 {
		requiresNode := root.CreateElement("Provides")
		for _, p := range jar.Provides {
			requiresNode.AddChild(util.BuildNodeWithText("Class", p))
		}
	}

	return root
}
