package analyzer

import (
	"binary"
	"fmt"
	"github.com/gnewton/jargo"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

const ANALYZERS_PATH = "C:\\Users\\Admin\\Work\\binfo\\backend"

func Jar(pathToJar string) *binary.JarBinary {
	jargoResult, err := jargo.GetJarInfo(pathToJar)
	if err != nil {
		log.Fatal(err)
	}

	jar := &binary.JarBinary{}
	jar.Architecture = ""
	jar.Filename = pathToJar
	jar.Dependencies = []binary.Dependency{}
	jar.Flags = []binary.Flag{}
	jar.Sections = []binary.Section{}
	jar.BuildBy = (*jargoResult.Manifest)["Build-By"]

	for index, file := range jargoResult.Files {
		fmt.Println(index, file)
	}

	//fmt.Printf("%+v\n", jar)
	return jar
}

func Dll(pathToDll string) *binary.PEBinary {
	return processWindowsBinary(pathToDll)
}

func Exe(pathToExe string) *binary.PEBinary {
	return processWindowsBinary(pathToExe)
}

func objdump(binaryFilePath string, args ...string) []byte {
	flagsString := "-" + strings.Join(args, "")
	command := createAnalyzerRunCommand("\\binutils\\objdump.exe " + flagsString + " " + binaryFilePath)
	cmd := exec.Command("cmd", "/C", command)
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	return stdoutStderr
}

func processWindowsBinary(pathToBinary string) *binary.PEBinary {
	bin := &binary.PEBinary{}
	byteDump := objdump(pathToBinary, "a", "f", "x")
	stringDump := string(byteDump)

	bin.Filename = pathToBinary
	bin.Dependencies = getDllDependencies(stringDump)
	bin.Architecture = getArchitecture(stringDump)
	bin.Sections = getSections(stringDump)
	bin.Flags = getFlags(stringDump)
	return bin
}

func createAnalyzerRunCommand(analyzer string) string {
	return "call " + ANALYZERS_PATH + analyzer
}

func getDllDependencies(dump string) []binary.Dependency {
	regex, _ := regexp.Compile("DLL Name: (.+?\\.dll)")
	matches := regex.FindAllStringSubmatch(dump, -1)
	dependencies := make([]binary.Dependency, len(matches))

	for index, element := range matches {
		dependencies[index] = binary.Dependency{Name: element[1]}
	}

	return dependencies
}

func getArchitecture(dump string) string {
	regex, _ := regexp.Compile("architecture: (.+?),")
	match := regex.FindStringSubmatch(dump)
	return match[1]
}

func getFlags(dump string) []binary.Flag {
	regex, _ := regexp.Compile("flags 0x[0-9a-f]+?:\\s+(([A-Z0-9_]+, )*[A-Z0-9_]+)\\s")
	flagsMatch := regex.FindStringSubmatch(dump)
	flagStrings := strings.Split(flagsMatch[1], ", ")
	flags := make([]binary.Flag, len(flagStrings))

	for index, element := range flagStrings {
		flags[index] = binary.Flag{Name: element}
	}

	return flags
}

func getSections(dump string) []binary.Section {
	return nil
}
