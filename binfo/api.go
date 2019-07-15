package binfo

import (
	"fmt"
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/util"
	"log"
	"runtime"
	"strings"
	"sync"
)

func ProcessFiles(binariesPath string, outputPath string) {
	bia := CreateAnalyzer(binariesPath, outputPath)
	done := make(chan bool)
	go bia.InitTasks()
	go bia.ProcessResults(done)
	var waitGroup sync.WaitGroup
	for i := 0; i < workersCount(); i++ {
		waitGroup.Add(1)
		go bia.fileProcessor(&waitGroup)
	}
	waitGroup.Wait()
	close(bia.Results)
	<-done
	fmt.Println("Done")
}

func (bia *BinaryInfoAnalyzer) InitTasks() {
	for _, binaryPath := range util.GetDirectoryFilePaths(bia.BinariesDirectory) {
		task := BinaryInfoTask{binaryPath}
		bia.Tasks <- task
	}
	close(bia.Tasks)
}

func (bia *BinaryInfoAnalyzer) ProcessResults(done chan bool) {
	for result := range bia.Results {
		if nil != result.Info {
			os.SaveResult(result.Info, bia.CreateOutputPath(result.FilePath))
		} else {
			log.Println("Cannot analyze file " + result.FilePath)
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

func createResult(filePath string, executable executable.Executable, err error) BinaryInfoResult {
	if nil != err {
		return BinaryInfoResult{FilePath: filePath}
	}
	return BinaryInfoResult{filePath, executable}
}
func workersCount() int {
	return runtime.NumCPU() + 1
}
