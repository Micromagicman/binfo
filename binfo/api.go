package binfo

import (
	"github.com/micromagicman/binary-info/analyzer"
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/logger"
	"github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/util"
	"runtime"
	"sync"
)

func ProcessFiles(binariesPath string, outputPath string) {
	os.CreateTemplateDirectory()
	os.InitOutputDirectory(outputPath)
	bia := CreateAnalyzer(binariesPath, outputPath)
	done := make(chan bool)
	// извлечь все внутренние jar-ники в директорию с бинарниками
	innerJars := analyzer.FindAndCreateInnerJars(binariesPath)
	go bia.InitTasks()
	go bia.ProcessResults(done)
	var waitGroup sync.WaitGroup
	workersCount := workersCount()
	for i := 0; i < workersCount; i++ {
		waitGroup.Add(1)
		go bia.fileProcessor(&waitGroup)
	}
	waitGroup.Wait()
	close(bia.Results)
	<-done
	// удалить все внутренние jar-ники из директории с бинарниками
	os.DeleteFiles(innerJars...)
	os.DeleteTemplateDirectory()
	logger.Info("Done")
}

func (bia *BinaryInfoAnalyzer) InitTasks() {
	for _, binaryPath := range util.GetDirectoryFilePaths(bia.BinariesDirectory) {
		task := BinaryInfoTask{binaryPath, nil}
		bia.Tasks <- task
	}
	close(bia.Tasks)
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