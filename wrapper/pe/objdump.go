package pe

import (
	"binfo/executable"
	"binfo/os"
	"binfo/wrapper"
	"github.com/decomp/exp/bin"
	"strconv"
	"strings"
)

type ObjDump struct {
	wrapper.BaseDump
}

func (od *ObjDump) GetName() string {
	return "objdump"
}

func (od *ObjDump) LoadFile(pathToExecutable string) bool {
	command := os.Exec.ObjDumpCommand(pathToExecutable, "-x")
	stdOut, err := os.Exec.Execute(command)
	if err != nil {
		return false
	}
	od.Content = string(stdOut)
	return true
}

func (od *ObjDump) Process(e executable.Executable) {
	peFile := e.(*executable.PortableExecutable)
	peFile.Architecture = od.getArchitecture()
	peFile.Exports = od.getExports()
	peFile.Flags = od.getFlags()
}

func (od *ObjDump) getArchitecture() string {
	return wrapper.Group(od.BaseDump.Find("architecture: (.+?),"), 1)
}

func (od *ObjDump) getFlags() []executable.Flag {
	flagsMatch := od.FindAll("flags 0x[0-9a-f]+?:\\s+(([A-Z0-9_]+, )*[A-Z0-9_]+)\\s")
	if len(flagsMatch) == 0 {
		return []executable.Flag{}
	}

	flagStrings := strings.Split(wrapper.Group(flagsMatch[0], 1), ", ")
	flags := make([]executable.Flag, len(flagStrings))

	for index, element := range flagStrings {
		flags[index] = executable.Flag{Name: element}
	}

	return flags
}


func (od *ObjDump) getExports() map[bin.Address]string {
	addresses := map[string]string{}
	names := map[string]string{}
	exports := map[bin.Address]string{}

	startIndex := strings.Index(od.Content, "Export Address Table")
	endIndex := strings.Index(od.Content, "PE File Base Relocations")
	dump := od.SubDump(startIndex, endIndex)

	addressesMatch := dump.FindAll(`\[\s*(\d+)\]\s\+base\[\s*\d+\]\s([a-f0-9]+)\sExport\sRVA`)
	namesMatch := dump.FindAll(`\n\s*\[\s*(\d+)\]\s([_a-zA-Z][^\s]+)`)

	for _, nm := range addressesMatch {
		addresses[wrapper.Group(nm, 1)] = wrapper.Group(nm, 2)
	}
	for _, em := range namesMatch {
		names[wrapper.Group(em, 1)] = wrapper.Group(em, 2)
	}

	for k, v := range addresses {
		if name, ok := names[k]; ok {
			address, _ := strconv.ParseInt(v, 16, 64)
			exports[bin.Address(address)] = name
		}
	}

	return exports
}
