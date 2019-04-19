package main

import (
	"analyzer"
	"binary"
	"log"
	"os"
	"path/filepath"
	"xml"
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) == 0 {
		log.Fatal("First argument must be a path to binary file")
		return
	}

	a := analyzer.CreateAnalyzer()
	a.CreateTemplateDirectory()

	for _, path := range arguments {
		if !checkFileExists(path) {
			log.Println("Cannot find binary file: " + path)
			continue
		}

		binaryPath := path
		binaryType := detectBinaryType(binaryPath)

		var binFile binary.Binary
		var err error
		if binaryType != analyzer.TYPE_UNKNOWN {
			binFile, err = a.Analyze(binaryPath, binaryType)
		} else {
			// TODO - попытка проанализировать файл без расширения
			log.Println("Unknown binary type for file " + binaryPath)
			binFile, err = a.TryToAnalyze(binaryPath)
		}

		if err != nil {
			log.Fatal("Cannot analyze file " + binaryPath + ": " + err.Error())
		}
		xml.BuildXml(binFile, filepath.Base(path) + ".xml")
	}
}

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func detectBinaryType(binaryPath string) int {
	switch filepath.Ext(binaryPath) {
		case ".dll": return analyzer.TYPE_DLL
		case ".lib": return analyzer.TYPE_LIB
		case ".jar": return analyzer.TYPE_JAR
		case ".exe": return analyzer.TYPE_EXE
		case ".ocx": return analyzer.TYPE_OCX
		case ".efi": return analyzer.TYPE_EFI
		case ".sys": return analyzer.TYPE_SYS
		case ".src": return analyzer.TYPE_SCR
		case ".drv": return analyzer.TYPE_DRV
		case ".cpl": return analyzer.TYPE_CPL
		case ".axf": return analyzer.TYPE_AXF
		case ".elf": return analyzer.TYPE_ELF
		case ".bin": return analyzer.TYPE_BIN
		case ".so": return analyzer.TYPE_SO
		case ".o": return analyzer.TYPE_O
		case ".a": return analyzer.TYPE_A
		default: return analyzer.TYPE_UNKNOWN
	}
}
