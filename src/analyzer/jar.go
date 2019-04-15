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

func JarAnalyzer(pathToJar string) *etree.Element {
	jarAnalyzerPath := ANALYZERS_PATH + "\\jaranalyzer"
	dir := filepath.Dir(pathToJar)
	ExecuteIn(jarAnalyzerPath + "\\runxmlsummary.bat " + dir + " " + TEMPLATE_DIRECTORY + "\\temp.xml", jarAnalyzerPath)
	return getJarFileElement(TEMPLATE_DIRECTORY + "\\temp.xml", pathToJar)
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
