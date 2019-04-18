package dump

import "binary"

type PEDump struct {
	BaseDump
}

func (pd *PEDump) GetSignature() string {
	signatureMatch := pd.Find("Signature:\\s+([^\\n]+)")
	return signatureMatch[1]
}

func (pd *PEDump) GetImportedFunctions() []binary.Function {
	functionsMatch := pd.FindAll("Function Name \\(Hint\\):\\s+([^\\s]+)")
	functions := make([]binary.Function, len(functionsMatch))

	for index, element := range functionsMatch {
		functions[index] = binary.Function{Name: element[1]}
	}

	return functions
}

func (pd *PEDump) GetTimestamp() int64 {
	return GetInteger(pd, "Timestamp:\\s+(\\d+)")
}

func (pd *PEDump) GetSize() int64 {
	return GetInteger(pd, "File size:\\s+(\\d+?) bytes")
}

func (pd *PEDump) getSize() int64 {
	return GetInteger(pd, "File size:\\s+(\\d+?) bytes")
}
