package dump

import (
	"binfo/binary"
)

type PEDump struct {
	BaseDump
}

func (pd *PEDump) GetImportedFunctions() []binary.Function {
	return pd.functionsByRegex("Function Name \\(Hint\\):\\s+([^\\s]+)")
}

func (pd *PEDump) GetExportedFunctions() []binary.Function {
	return pd.functionsByRegex("Function Name:\\s+([^\\s]+)")
}

func (pd *PEDump) GetTimestamp() int64 {
	return GetInteger(pd, "Timestamp:\\s+(\\d+)")
}

func (pd *PEDump) GetSize() int64 {
	return GetInteger(pd, "File size:\\s+(\\d+?) bytes")
}

func (pd *PEDump) GetEntryPointAddress() string {
	return Group(pd.Find("Address of entry point:\\s+(0x[^\\n]+)"), 1)
}

func (pd *PEDump) GetCodeSectionAddress() string {
	return Group(pd.Find("Base address of code section:\\s+(0x[^\\n]+)"), 1)
}

func (pd *PEDump) GetDataSectionAddress() string {
	return Group(pd.Find("Base address of data section:\\s+(0x[^\\n]+)"), 1)
}

func (pd *PEDump) functionsByRegex(regex string) []binary.Function {
	functionsMatch := pd.FindAll(regex)
	functions := make([]binary.Function, len(functionsMatch))

	for index, element := range functionsMatch {
		functions[index] = binary.Function{Name: Group(element, 1)}
	}

	return functions
}