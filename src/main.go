package main

import (
	"analyzer"
	"binary"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"util"
)

func main() {
	binaryDirectory := flag.String("d", "", "Directory with binary files")
	outputDirectory := flag.String("o", "", "Output directory")
	flag.Parse()

	if "" == *binaryDirectory {
		log.Fatal("Directory with binary files not defined")
	}

	if "" == *outputDirectory {
		log.Fatal("Output directory not defined")
	}

	binariesPath := *binaryDirectory
	outputPath := *outputDirectory
	a := analyzer.CreateAnalyzer()
	a.CreateTemplateDirectory()
	a.InitOutputDirectory(outputPath)

	for _, path := range util.GetDirectoryFilePaths(binariesPath) {
		if !util.CheckFileExists(path) {
			log.Println("Cannot find binary file: " + path)
			continue
		}

		fmt.Println("Processing " + path)
		var binFile binary.Binary
		var err error
		binaryType := detectBinaryType(path)

		if binaryType != analyzer.TYPE_UNKNOWN {
			binFile, err = a.Analyze(path, binaryType)
		} else {
			// попытка проанализировать файл без расширения
			log.Println("Unknown binary type for file " + path)
			binFile, err = a.TryToAnalyze(path)
		}

		util.LogIfError(err, "Cannot analyze file " +path)
		a.SaveResult(binFile, path)
	}

	a.DeleteTemplateDirectory()
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
