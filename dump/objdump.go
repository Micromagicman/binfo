package dump

import (
	"binfo/binary"
	"strings"
)

type ObjDump struct {
	BaseDump
}

func (od *ObjDump) GetDependencies() []binary.Dependency {
	depMatches := od.BaseDump.FindAll("DLL Name: (.+?\\.dll)")
	dependencies := make([]binary.Dependency, len(depMatches))

	for index, element := range depMatches {
		dependencies[index] = binary.Dependency{Name: Group(element,1)}
	}

	return dependencies
}

func (od *ObjDump) GetArchitecture() string {
	return Group(od.BaseDump.Find("architecture: (.+?),"), 1)
}

func (od *ObjDump) GetFlags() []binary.Flag {
	flagsMatch := od.FindAll("flags 0x[0-9a-f]+?:\\s+(([A-Z0-9_]+, )*[A-Z0-9_]+)\\s")
	if len(flagsMatch) == 0 {
		return []binary.Flag{}
	}
	
	flagStrings := strings.Split(Group(flagsMatch[0], 1), ", ")
	flags := make([]binary.Flag, len(flagStrings))

	for index, element := range flagStrings {
		flags[index] = binary.Flag{Name: element}
	}

	return flags
}
