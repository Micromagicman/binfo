package main

import (
	"analyzer"
	"binary"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	//exeFile := analyzer.Exe("C:\\Users\\Admin\\Work\\binfo\\backend\\binutils\\objdump.exe")
	dllFile := analyzer.Dll("C:\\Users\\Admin\\Work\\binfo\\backend\\binutils\\libiconv-2.dll")
	fmt.Println(dllFile)
	createXml(dllFile)
	//fmt.Println(analyzer.Dll("C:\\Users\\Admin\\Work\\binfo\\backend\\binutils\\libiconv-2.dll"))
}

func createXml(bin *binary.BinaryFile) {
	file, _ := os.Create("test.xml")
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	root := doc.CreateElement("Binary")
	filenameElem := root.CreateElement("Filename")
	filenameElem.CreateText(bin.Filename)
	architectureElem := root.CreateElement("Architecture")
	architectureElem.CreateText(bin.Architecture)

	flagsElem := root.CreateElement("Flags")
	for index, flag := range bin.Flags {
		flagElem := flagsElem.CreateElement("Flag")
		flagElem.CreateAttr("id", strconv.Itoa(index))
		flagElem.CreateText(flag.Name)
	}

	dependenciesElem := root.CreateElement("Dependencies")
	for index, dependency := range bin.Dependencies {
		dependencyElem := dependenciesElem.CreateElement("Dependency")
		dependencyElem.CreateAttr("id", strconv.Itoa(index))
		dependencyElem.CreateText(dependency.Name)
	}

	doc.Indent(4)
	_, err := doc.WriteTo(file)
	if err != nil {
		log.Fatal("Error when create xml")
	}
}



