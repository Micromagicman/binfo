package analyzer

import (
	"fmt"
	"github.com/beevik/etree"
	"path/filepath"
	"strings"
)

func (a *Analyzer) JarAnalyzer(pathToJar string) (*etree.Element, error) {
	jarAnalyzerPath := a.Executor.AnalyzersPath + "jaranalyzer\\"
	dir := filepath.Dir(pathToJar)
	_, executeError := a.Executor.Execute(jarAnalyzerPath+"runxmlsummary.bat "+dir+" "+a.Executor.TemplateDirectory + "temp.xml")
	fmt.Println(jarAnalyzerPath+"runxmlsummary.bat "+dir+" "+a.Executor.TemplateDirectory + "temp.xml")

	if executeError != nil {
		fmt.Println(executeError.Error())
		return nil, executeError
	}

	return getJarFileElement(a.Executor.TemplateDirectory + "temp.xml", pathToJar)
}

func getJarFileElement(pathToJarAnalyzerXml string, pathToJar string) (*etree.Element, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(pathToJarAnalyzerXml); err != nil {
		return nil, err
	}

	for _, jar := range doc.FindElements("//Jar") {
		if strings.HasSuffix(pathToJar, jar.SelectAttr("name").Value) {
			return jar.ChildElements()[0], nil // Summary
		}
	}

	return nil, nil
}
