package common

import (
	"binfo/executable"
	"os"
	"path/filepath"
)

type FileStat struct {
	Filename string
	FileInfo os.FileInfo
}

func (fs *FileStat) GetName() string {
	return "filestat"
}

func (fs *FileStat) LoadFile(pathToExecutable string) bool {
	file, err := os.Open(pathToExecutable)
	if err != nil {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	fs.Filename, _ = filepath.Abs(pathToExecutable)
	fs.FileInfo = info
	return true
}

func (fs *FileStat) Process(e executable.Executable) {
	e.SetFileName(fs.Filename)
	e.SetSize(fs.FileInfo.Size())
}