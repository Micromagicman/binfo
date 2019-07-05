package pe

import (
	"github.com/mewrev/pe"
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/util"
)

type MemrevPE struct {
	File *pe.File
}

func (mpe *MemrevPE) GetName() string {
	return "memrev/pe"
}

func (mpe *MemrevPE) LoadFile(pathToExecutable string) bool {
	file, err := pe.Open(pathToExecutable)
	if nil != err {
		return false
	}
	mpe.File = file
	return true
}

func (mpe *MemrevPE) Process(e executable.Executable) {
	peFile := e.(*executable.PortableExecutable)
	fileHeader, _ := mpe.File.FileHeader()
	peFile.SectionNumber = fileHeader.NSection
	peFile.Architecture = fileHeader.Arch.String()

	sectionHeaders, _ := mpe.File.SectHeaders()
	optionalHeader, _ := mpe.File.OptHeader()

	peFile.LinkerVersion = util.Int64ToString(int64(optionalHeader.MajorLinkVer)) + "." + util.Int64ToString(int64(optionalHeader.MinorLinkVer))
	peFile.OsVersion = util.Int64ToString(int64(optionalHeader.MajorOSVer)) + "." + util.Int64ToString(int64(optionalHeader.MinorOSVer))
	peFile.Checksum = util.Int64ToHex(int64(optionalHeader.Checksum))
	peFile.CodeRVA = util.Int64ToHex(int64(optionalHeader.CodeBase))
	peFile.CodeSize = util.Int64ToString(int64(optionalHeader.CodeSize))
	peFile.DataRVA = util.Int64ToHex(int64(optionalHeader.DataSize))
	peFile.DataSize = util.Int64ToString(int64(optionalHeader.DataSize))

	peFile.Sections = sectionHeaders
}
