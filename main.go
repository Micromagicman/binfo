package main

import (
	"binfo/analyzer"
	"binfo/executable"
	"binfo/os"
	"binfo/util"
	"flag"
	"fmt"
	"log"
	"path/filepath"
)

func main() {
	binaryDirectory := flag.String("d", "", "Directory with executable files")
	outputDirectory := flag.String("o", "", "Output directory")
	flag.Usage = helpMessage
	flag.Parse()

	if "" == *binaryDirectory || "" == *outputDirectory {
		flag.Usage()
		return
	}

	binariesPath := *binaryDirectory
	outputPath := *outputDirectory
	a := analyzer.CreateAnalyzer()
	os.CreateTemplateDirectory()
	os.InitOutputDirectory(outputPath)

	for _, path := range util.GetDirectoryFilePaths(binariesPath) {
		if !util.CheckFileExists(path) {
			log.Println("Cannot find executable file: " + path)
			continue
		}

		fmt.Println("Processing " + path)
		var binFile executable.Executable
		var err error
		binaryType := detectBinaryType(path)

		if binaryType != analyzer.TYPE_UNKNOWN {
			binFile, err = a.Analyze(path, binaryType)
		} else {
			// попытка проанализировать файл без расширения
			log.Println("Unknown executable type for file " + path)
			binFile, err = a.TryToAnalyze(path)
			if nil == binFile {
				continue
			}
		}

		util.LogIfError(err, "Cannot analyze file "+path)
		os.SaveResult(binFile, outputPath, path)
	}

	os.DeleteTemplateDirectory()
}

func detectBinaryType(binaryPath string) int {
	switch filepath.Ext(binaryPath) {
	case ".dll":
		return analyzer.TYPE_DLL
	case ".lib":
		return analyzer.TYPE_LIB
	case ".jar":
		return analyzer.TYPE_JAR
	case ".exe":
		return analyzer.TYPE_EXE
	case ".ocx":
		return analyzer.TYPE_OCX
	case ".efi":
		return analyzer.TYPE_EFI
	case ".sys":
		return analyzer.TYPE_SYS
	case ".src":
		return analyzer.TYPE_SCR
	case ".drv":
		return analyzer.TYPE_DRV
	case ".cpl":
		return analyzer.TYPE_CPL
	case ".axf":
		return analyzer.TYPE_AXF
	case ".elf":
		return analyzer.TYPE_ELF
	case ".bin":
		return analyzer.TYPE_BIN
	case ".so":
		return analyzer.TYPE_SO
	case ".o":
		return analyzer.TYPE_O
	case ".a":
		return analyzer.TYPE_A
	default:
		return analyzer.TYPE_UNKNOWN
	}
}

func helpMessage() {
	fmt.Println(
			"Usage:\n" +
				"\tmain.exe -d path/to/executable -o path/to/output\n" +
		"Flags:\n" +
			"\t-d directory with executables (required)\n" +
			"\t-o output directory (required)\n" +
			"Supported Formats:\n" +
			"\t- Windows Portable Executable: exe, dll, ocx, sys, scr, drv, cpl, efi\n" +
			"\t- Executable Linkable: exe, so, axf, bin, elf, o, a, prx\n" +
			"\t- jar")
}
