package main

import (
	"analyzer"
	"xml"
)

const (
	TYPE_EXE = 0
	TYPE_DLL = 1
	TYPE_JAR = 2
)

func main() {
	//exeFile := analyzer.Exe("C:\\Users\\Admin\\Work\\binfo\\backend\\binutils\\objdump.exe")
	analyze( ,TYPE_DLL)
	dllFile := analyzer.Dll("C:\\Users\\Admin\\Work\\binfo\\backend\\binutils\\libiconv-2.dll")
	//jarFile := analyzer.Jar("C:\\Users\\Admin\\Work\\jars\\drone-java.jar")
	xml.BuildXml(dllFile, "output_dll.xml")
	//xml.BuildXml(jarFile, "output_jar.xml")
	//analyzer.Jar("C:\\Users\\Admin\\Work(\\jars\\drone-java.jar")
	//fmt.Println(analyzer.Dll("C:\\Users\\Admin\\Work\\binfo\\backend\\binutils\\libiconv-2.dll"))
}

func detectAnalyzerType(filename string) {

}

func analyze(pathToBinary string, binaryType int) {
	if binaryType == TYPE_EXE {
		file := analyzer.Exe(pathToBinary)
		xml.BuildXml(file, "output.xml")
	} else if binaryType == TYPE_DLL {
		file := analyzer.Dll(pathToBinary)
		xml.BuildXml(file, "output.xml")
	} else if binaryType == TYPE_JAR {
		// TODO
	}
}

