package binfo

import (
	"github.com/micromagicman/binary-info/analyzer"
	"github.com/micromagicman/binary-info/logger"
	"github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/util"
	"strings"
	"sync"
)

func CreateAnalyzer(binariesDirectory string, outputDirectory string) *BinaryInfoAnalyzer {
	return &BinaryInfoAnalyzer{
		binariesDirectory,
		outputDirectory,
		analyzer.CreateAnalyzer(),
		make(chan BinaryInfoTask, 100),
		make(chan BinaryInfoResult, 100),
	}
}

func (bia *BinaryInfoAnalyzer) ProcessResults(done chan bool) {
	for result := range bia.Results {
		if nil != result.Info {
			os.SaveResult(result.Info, bia.CreateOutputPath(result.FilePath))
		} else {
			logger.Warning("Cannot analyze file " + result.FilePath)
		}
	}
	done <- true
}

func (bia *BinaryInfoAnalyzer) CreateOutputPath(binaryPath string) string {
	fileName := strings.TrimPrefix(binaryPath, bia.BinariesDirectory)
	return util.CreateOutputPath(fileName, bia.OutputDirectory)
}

func (bia *BinaryInfoAnalyzer) fileProcessor(waitGroup *sync.WaitGroup) {
	for task := range bia.Tasks {
		binFile, err :=  bia.FileAnalyzer.Analyze(task.FilePath)
		bia.Results <- createResult(task.FilePath, binFile, err)
	}
	waitGroup.Done()
}
