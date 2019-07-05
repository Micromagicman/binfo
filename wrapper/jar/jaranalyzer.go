package jar

import (
	"github.com/beevik/etree"
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/wrapper"
	"path/filepath"
	"strings"
)

type JarAnalyzer struct {
	wrapper.OnlyRun
	CurrentTree *etree.Element
	JarElements []*etree.Element
}

func (ja *JarAnalyzer) GetWindowsCommand() string {
	return "call " + os.BackendDir + os.Sep + "jaranalyzer" + os.Sep + "runxmlsummary.bat"
}

func (ja *JarAnalyzer) GetLinuxCommand() string {
	return "call " + os.BackendDir + os.Sep + "jaranalyzer" + os.Sep + "runxmlsummary.sh"
}

func (ja *JarAnalyzer) GetName() string {
	return "jaranalyzer"
}

func (ja *JarAnalyzer) LoadFile(pathToExecutable string) bool {
	arguments := []string{
		filepath.Dir(pathToExecutable),
		os.TemplateDir + os.Sep + "temp.xml",
	}
	if !ja.WasExecuted() {
		if _, err := os.Execute(ja, arguments...); nil != err {
			return false
		}
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(os.TemplateDir + os.Sep + "temp.xml"); nil != err {
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