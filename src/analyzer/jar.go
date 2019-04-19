package analyzer

import (
	"github.com/beevik/etree"
	"github.com/gnewton/jargo"
	"log"
	"path/filepath"
	"strings"
)

func Jargo(pathToJar string) *jargo.JarInfo {
	jargoResult, err := jargo.GetJarInfo(pathToJar)
	if err != nil {
		log.Fatal(err)
	}
	return jargoResult
}

func (a *Analyzer) JarAnalyzer(pathToJar string) (*etree.Element, error) {
	jarAnalyzerPath := a.Executor.AnalyzersPath + a.Executor.Sep + "jaranalyzer"
	dir := filepath.Dir(pathToJar)
	_, executeError := a.Executor.ExecuteIn(jarAnalyzerPath+a.Executor.Sep+"runxmlsummary.bat "+dir+" "+a.Executor.TemplateDirectory+a.Executor.Sep+"temp.xml", jarAnalyzerPath)

	if executeError != nil {
		return nil, executeError
	}

	return getJarFileElement(a.Executor.TemplateDirectory + a.Executor.Sep + "temp.xml", pathToJar), nil
}

func getJarFileElement(pathToJarAnalyzerXml string, pathToJar string) *etree.Element {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(pathToJarAnalyzerXml); err != nil {
		panic(err)
	}

	for _, jar := range doc.FindElements("//Jar") {
		if strings.HasSuffix(pathToJar, jar.SelectAttr("name").Value) {
			return jar.ChildElements()[0] // Summary
		}
	}

	return nil
}
