package analyzer

import (
	"binfo/dump"
	"log"
	"path/filepath"
	"strings"
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

func (a *Analyzer) Tattletale(jarFilePath string) *dump.Tattletale {
	command := a.Executor.TattletaleCommand(jarFilePath)
	_, _ = a.Executor.Execute(command)
	jarHtmlReport := "jar" + a.Executor.Sep + filepath.Base(jarFilePath) + ".html"
	wrapper, err := dump.CreateTattletaleWrapper(a.Executor.TemplateDirectory + jarHtmlReport)

	if err != nil {
		log.Println("Cannot analyze " + jarFilePath + " via Tattletale")
		return nil
	}

	return wrapper
}