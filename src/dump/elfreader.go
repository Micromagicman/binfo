package dump

type ELFReader struct {
	BaseDump
}

func (ed *ELFReader) GetOperatingSystem() string {
	return Group(ed.Find("Operating System:\\s+([^.]+)"), 1)
}

func (ed *ELFReader) GetFormat() string {
	return Group(ed.Find("Bit format:\\s+([^.]+)"), 1)
}

func (ed *ELFReader) GetEndianess() string {
	return Group(ed.Find("Endianess:\\s+([^.]+)"), 1)
}

func (ed *ELFReader) GetVersion() string {
	return Group(ed.Find("ELF Version:\\s+([^.]+)"), 1)
}

func (ed *ELFReader) GetType() string {
	return Group(ed.Find("ELF Type:\\s+([^.]+)"), 1)
}
