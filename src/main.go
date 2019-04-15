package main

import (
	"analyzer"
	"binary"
	"xml"
)

const (
	TYPE_EXE = 0
	TYPE_DLL = 1
	TYPE_JAR = 2
)

func main() {
	analyzer.CreateTemplateDirectory()
	analyze("C:\\Users\\Admin\\Work\\jars\\drone-java.jar", TYPE_JAR)
}

func analyze(pathToBinary string, binaryType int) {
	var file binary.Binary
	if binaryType == TYPE_EXE {
		file = analyzer.Exe(pathToBinary)
	} else if binaryType == TYPE_DLL {
		file = analyzer.Dll(pathToBinary)
	} else if binaryType == TYPE_JAR {
		file = analyzer.Jar(pathToBinary)
	}

	if file != nil {
		xml.BuildXml(file, "output.xml")
	}
}

