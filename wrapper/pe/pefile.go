package pe

import (
	"github.com/H5eye/go-pefile"
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/util"
)

type PEFile struct {
	PEFile *pefile.PE
}

func (pf *PEFile) GetName() string {
	return "H5eye/go-pefile"
}

func (pf *PEFile) LoadFile(pathToExecutable string) bool {
	peLoad, err := pefile.Load(pathToExecutable)
	if nil != err {
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
