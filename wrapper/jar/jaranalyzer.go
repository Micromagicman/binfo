package jar

import (
	"binfo/executable"
	"binfo/os"
	"binfo/wrapper"
	"fmt"
	"github.com/beevik/etree"
	"path/filepath"
	"strings"
)

type JarAnalyzer struct {
	wrapper.OnlyRun
	CurrentTree *etree.Element
	JarElements []*etree.Element
}

func (ja *JarAnalyzer) GetWindowsCommand(filePath string) string {
	fileDir := filepath.Dir(filePath)
	fmt.Println("call " + os.BackendDir + os.Sep + "jaranalyzer" + os.Sep + "runxmlsummary.bat " +
		fileDir + " " + os.TemplateDir)
	return "call " + os.BackendDir + os.Sep + "jaranalyzer" + os.Sep + "runxmlsummary.bat " +
		fileDir + " " + os.TemplateDir + os.Sep + "temp.xml"
}

func (ja *JarAnalyzer) GetLinuxCommand(filePath string) string {
	fileDir := filepath.Dir(filePath)
	return "call " + os.BackendDir + os.Sep + "jaranalyzer" + os.Sep + "runxmlsummary.bat "+
		fileDir + " " + os.TemplateDir + os.Sep + "temp.xml"
}

func (ja *JarAnalyzer) GetName() string {
	return "jaranalyzer"
}

func (ja *JarAnalyzer) LoadFile(pathToExecutable string) bool {
	if !ja.WasExecuted() {
		_, err := os.Execute(pathToExecutable, ja)
		if err != nil {
			return false
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromFile(os.TemplateDir + os.Sep + "temp.xml"); err != nil {
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