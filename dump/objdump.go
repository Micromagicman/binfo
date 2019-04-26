package dump

import (
	"binfo/executable"
	"strings"
)

type ObjDump struct {
	BaseDump
}

func (od *ObjDump) GetDependencies() []executable.Dependency {
	depMatches := od.BaseDump.FindAll("DLL Name: (.+?\\.dll)")
	dependencies := make([]executable.Dependency, len(depMatches))

	for index, element := range depMatches {
		dependencies[index] = executable.Dependency{Name: Group(element,1)}
	}

	return dependencies
}

func (od *ObjDump) GetArchitecture() string {
	return Group(od.BaseDump.Find("architecture: (.+?),"), 1)
}

func (od *ObjDump) GetFlags() []executable.Flag {
	flagsMatch := od.FindAll("flags 0x[0-9a-f]+?:\\s+(([A-Z0-9_]+, )*[A-Z0-9_]+)\\s")
	if len(flagsMatch) == 0 {
		return []executable.Flag{}
	}
	
	flagStrings := strings.Split(Group(flagsMatch[0], 1), ", ")
	flags := make([]executable.Flag, len(flagStrings))

	for index, element := range flagStrings {
		flags[index] = executable.Flag{Name: element}
	}

	return flags
}
