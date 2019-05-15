package wrapper

import (
	"binfo/executable"
	"binfo/util"
	"github.com/mewrev/pe"
)

type MemrevPE struct {
	File *pe.File
}

func CreateMemrevPEWrapper(pathToBinary string) (*MemrevPE, error) {
	file, err := pe.Open(pathToBinary)
	if err != nil {
		return nil, err
	}
	return &MemrevPE{file}, nil
}

func (pe *MemrevPE) Process(peBin *executable.PortableExecutable) {
	fileHeader, _ := pe.File.FileHeader()
	peBin.SectionNumber = fileHeader.NSection
	peBin.Architecture = fileHeader.Arch.String()

	sectionHeaders, _ := pe.File.SectHeaders()
	optionalHeader, _ := pe.File.OptHeader()

	peBin.LinkerVersion = util.Int64ToString(int64(optionalHeader.MajorLinkVer)) + "." + util.Int64ToString(int64(optionalHeader.MinorLinkVer))
	peBin.OsVersion = util.Int64ToString(int64(optionalHeader.MajorOSVer)) + "." + util.Int64ToString(int64(optionalHeader.MinorOSVer))
	peBin.Checksum = util.Int64ToHex(int64(optionalHeader.Checksum))
	peBin.CodeRVA = util.Int64ToHex(int64(optionalHeader.CodeBase))
	peBin.CodeSize = util.Int64ToString(int64(optionalHeader.CodeSize))
	peBin.DataRVA = util.Int64ToHex(int64(optionalHeader.DataSize))
	peBin.DataSize = util.Int64ToString(int64(optionalHeader.DataSize))

	peBin.Sections = sectionHeaders
}
