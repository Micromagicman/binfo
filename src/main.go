package main

import (
	"analyzer"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"xml"
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) == 0 {
		log.Fatal("First argument must be a path to binary file")
		return
	}

	analyzer.CreateTemplateDirectory()
	for index, path := range arguments {
		binaryPath := path
		binaryType := detectBinaryType(binaryPath)
		if binaryType == analyzer.TYPE_UNKNOWN {
			// TODO - попытка проанализировать файл без расширения
			log.Fatal("Unknown binary type for file " + binaryPath)
		}

		bin := analyzer.Analyze(binaryPath, binaryType)
		if bin != nil {
			xml.BuildXml(bin, strconv.Itoa(index) + ".xml")
		}
	}

}

func detectBinaryType(binaryPath string) int {
	extension := filepath.Ext(binaryPath)
	switch extension {
		case ".dll": return analyzer.TYPE_DLL
		case ".jar": return analyzer.TYPE_JAR
		case ".exe": return analyzer.TYPE_EXE
	}
	return analyzer.TYPE_UNKNOWN
}
