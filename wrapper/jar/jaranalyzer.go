package jar

import (
	"binfo/executable"
	osUtils "binfo/os"
	"binfo/wrapper"
	"github.com/beevik/etree"
	"path/filepath"
	"strings"
)

type JarAnalyzer struct {
	wrapper.OnlyRun
	CurrentTree *etree.Element
	JarElements []*etree.Element
}

func (ja *JarAnalyzer) GetName() string {
	return "jaranalyzer"
}

func (ja *JarAnalyzer) LoadFile(pathToExecutable string) bool {
	if !ja.WasExecuted() {
		jarAnalyzerPath := osUtils.Exec.AnalyzersPath + "jaranalyzer" + osUtils.Exec.Sep
		dir := filepath.Dir(pathToExecutable)
		_, err := osUtils.Exec.Execute(jarAnalyzerPath + "runxmlsummary.bat " + dir + " " + osUtils.Exec.TemplateDirectory + "temp.xml")
		if err != nil {
			return false
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromFile(osUtils.Exec.TemplateDirectory + "temp.xml"); err != nil {
			return false
		}

		ja.JarElements = doc.FindElements("//Jar")
		ja.MarkAsExecuted()
	}

	for _, jar := range ja.JarElements {
		if strings.HasSuffix(pathToExecutable, jar.SelectAttr("name").Value) {
			ja.CurrentTree = jar.ChildElements()[0]
			return true
		}
	}

	return false
}

func (ja *JarAnalyzer) Process(e executable.Executable) {
	jarFile := e.(*executable.JarExecutable)
	jarFile.JarAnalyzerTree = ja.CurrentTree
}