package analyzer

import (
	"binfo/dump"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

func (a *Analyzer) ObjDump(binaryFilePath string, args ...string) *dump.ObjDump {
	flagsString := "-" + strings.Join(args, "")
	command := a.Executor.ObjDumpCommand(binaryFilePath, flagsString)
	stdOut, _ := a.Executor.Execute(command)

	objDump := &dump.ObjDump{}
	objDump.Content = string(stdOut)
	return objDump
}

func (a *Analyzer) PEDumper(binaryFilePath string) *dump.PEDump {
	command := a.Executor.PEDumperCommand(binaryFilePath)
	stdOut, _ := a.Executor.Execute(command)

	peDump := &dump.PEDump{}
	peDump.Content = string(stdOut)
	return peDump
}

func (a *Analyzer) ELFReader(binaryFilePath string) *dump.ELFReader {
	command := a.Executor.ELFReaderCommand(binaryFilePath)
	stdOut, _ := a.Executor.Execute(command)

	elfDump := &dump.ELFReader{}
	elfDump.Content = string(stdOut)
	return elfDump
}

func (a *Analyzer) CDetect(binaryFilePath string) string {
	command := a.Executor.CDetectCommand(binaryFilePath)
	stdOut, _ := a.Executor.Execute(command)
	return string(stdOut)
}

func (a *Analyzer) Tattletale(jarFilePath string) *dump.Tattletale {
	if !a.Cache.Tattletale {
		command := a.Executor.TattletaleCommand(jarFilePath)
		_, _ = a.Executor.Execute(command)
		a.Cache.Tattletale = true
	}
	jarHtmlReport := "jar" + a.Executor.Sep + filepath.Base(jarFilePath) + ".html"
	wrapper, err := dump.CreateTattletaleWrapper(a.Executor.TemplateDirectory + jarHtmlReport)

	if err != nil {
		log.Println("Cannot analyze " + jarFilePath + " via Tattletale")
		return nil
	}

	return wrapper
}
func (a *Analyzer) JarAnalyzer(pathToJar string) (*etree.Element, error) {
	if !a.Cache.JarAnalyzer {
		jarAnalyzerPath := a.Executor.AnalyzersPath + "jaranalyzer\\"
		dir := filepath.Dir(pathToJar)
		_, executeError := a.Executor.Execute(jarAnalyzerPath + "runxmlsummary.bat " + dir + " " + a.Executor.TemplateDirectory + "temp.xml")

		if executeError != nil {
			fmt.Println(executeError.Error())
			return nil, executeError
		}

		a.Cache.JarAnalyzer = true
	}

	return getJarFileElement(a.Executor.TemplateDirectory+"temp.xml", pathToJar)
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
