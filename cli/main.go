package main

import (
	"flag"
	"fmt"
	"github.com/micromagicman/binary-info/binfo"
)

func main() {
	binaryDirectory := flag.String("d", "", "Directory with executable files")
	outputDirectory := flag.String("o", "", "Output directory")
	flag.Usage = helpMessage
	flag.Parse()
	if "" != *binaryDirectory && "" != *outputDirectory {
		binfo.ProcessFiles(*binaryDirectory, *outputDirectory)
	} else {
		flag.Usage()
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
