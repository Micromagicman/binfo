package binfo

import (
	"github.com/micromagicman/binary-info/analyzer"
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/os"
)

type BinaryInfoTask struct {
	FilePath string
}

type BinaryInfoResult struct {
	FilePath string
	Info     executable.Executable
}

type BinaryInfoAnalyzer struct {
	BinariesDirectory string
	OutputDirectory   string
	FileAnalyzer      *analyzer.Analyzer
	Tasks             chan BinaryInfoTask
	Results           chan BinaryInfoResult
}

func CreateAnalyzer(binariesDirectory string, outputDirectory string) *BinaryInfoAnalyzer {
	os.CreateTemplateDirectory()
	os.InitOutputDirectory(outputDirectory)
	return &BinaryInfoAnalyzer{
		binariesDirectory,
		outputDirectory,
		analyzer.CreateAnalyzer(),
		make(chan BinaryInfoTask, 100),
		make(chan BinaryInfoResult, 100),
	}
}