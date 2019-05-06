package dump

import (
	"binfo/executable"
)

type PEDump struct {
	BaseDump
}

func (pd *PEDump) GetImportedFunctions() []executable.Function {
	return pd.functionsByRegex("Function Name \\(Hint\\):\\s+([^\\s]+)")
}

func (pd *PEDump) GetExportedFunctions() []executable.Function {
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

func (pd *PEDump) functionsByRegex(regex string) []executable.Function {
	functionsMatch := pd.FindAll(regex)
	functions := make([]executable.Function, len(functionsMatch))

	for index, element := range functionsMatch {
		functions[index] = executable.Function{Name: Group(element, 1)}
	}

	return functions
}
