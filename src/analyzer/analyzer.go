package analyzer

import (
	"binary"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

const ANALYZERS_PATH = "C:\\Users\\Admin\\Work\\binfo\\backend"

func Jar(pathToJar string) {

}

func Dll(pathToDll string) *binary.BinaryFile {
	return processWindowsBinary(pathToDll)
}

func Exe(pathToExe string) *binary.BinaryFile {
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

func processWindowsBinary(pathToBinary string) *binary.BinaryFile {
	bin := &binary.BinaryFile{}
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
