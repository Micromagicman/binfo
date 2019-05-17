package dump

import (
	"binfo/executable"
	"github.com/decomp/exp/bin"
	"strconv"
	"strings"
)

type ObjDump struct {
	BaseDump
}

func (od *ObjDump) GetDependencies() []string {
	depMatches := od.BaseDump.FindAll("DLL Name: (.+?\\.dll)")
	dependencies := make([]string, len(depMatches))

	for index, element := range depMatches {
		dependencies[index] = Group(element, 1)
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


func (od *ObjDump) GetExports() map[bin.Address]string {
	addresses := map[string]string{}
	names := map[string]string{}
	exports := map[bin.Address]string{}

	startIndex := strings.Index(od.Content, "Export Address Table")
	endIndex := strings.Index(od.Content, "PE File Base Relocations")
	dump := od.SubDump(startIndex, endIndex)

	addressesMatch := dump.FindAll(`\[\s*(\d+)\]\s\+base\[\s*\d+\]\s([a-f0-9]+)\sExport\sRVA`)
	namesMatch := dump.FindAll(`\n\s*\[\s*(\d+)\]\s([_a-zA-Z][^\s]+)`)

	for _, nm := range addressesMatch {
		addresses[Group(nm, 1)] = Group(nm, 2)
	}
	for _, em := range namesMatch {
		names[Group(em, 1)] = Group(em, 2)
	}

	for k, v := range addresses {
		if name, ok := names[k]; ok {
			address, _ := strconv.ParseInt(v, 16, 64)
			exports[bin.Address(address)] = name
		}
	}

	return exports
}