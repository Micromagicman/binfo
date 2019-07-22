package binfo

import (
	"github.com/micromagicman/binary-info/analyzer"
	"github.com/micromagicman/binary-info/executable"
)

type BinaryInfoTask struct {
	FilePath   string
	ChildTasks []*BinaryInfoTask
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
