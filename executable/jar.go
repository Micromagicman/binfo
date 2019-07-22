package executable

import (
	"github.com/micromagicman/binary-info/util"
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
	Children        []*JarExecutable
}

func (jar *JarExecutable) GetType() string {
	return "JAR Archive"
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
		xmlKey := strings.Replace(key, "-", "", -1)
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
		requiresNode := root.CreateElement("Imports")
		for _, r := range jar.Requires {
			requiresNode.AddChild(util.BuildNodeWithText("Class", r))
		}
	}

	if len(jar.Children) > 0 {
		innerJarsNode := root.CreateElement("InnerJars")
		for _, jc := range jar.Children {
			childDocument := etree.NewDocument()
			innerJarsNode.AddChild(jc.BuildXml(childDocument))
		}
 	}

	if len(jar.Provides) > 0 {
		requiresNode := root.CreateElement("Exports")
		for _, p := range jar.Provides {
			requiresNode.AddChild(util.BuildNodeWithText("Class", p))
		}
	}

	return root
}
