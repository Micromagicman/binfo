package pe

import (
	"binfo/executable"
	"binfo/util"
	"github.com/H5eye/go-pefile"
)

type PEFile struct {
	PEFile *pefile.PE
}

func (pf *PEFile) GetName() string {
	return "H5eye/go-pefile"
}

func (pf *PEFile) LoadFile(pathToExecutable string) bool {
	peLoad, err := pefile.Load(pathToExecutable)
	if err != nil {
		return false
	}
	pf.PEFile = peLoad
	return true
}

func (pf *PEFile) Process(e executable.Executable) {
	peFile := e.(*executable.PortableExecutable)
	peFile.Timestamp = int64(pf.PEFile.FileHeader.TimeDateStamp)
	peFile.Time = util.TimestampToTime(peFile.Timestamp)
}
