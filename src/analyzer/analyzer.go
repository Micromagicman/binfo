package analyzer

import (
	"binary"
	"fmt"
	"github.com/mewrev/pe"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TYPE_UNKNOWN = -1
	TYPE_EXE = 0
	TYPE_DLL = 1
	TYPE_OCX = 2
	TYPE_SYS = 3
	TYPE_SCR = 4
	TYPE_DRV = 5
	TYPE_CPL = 6
	TYPE_EFI = 7
	TYPE_JAR = 10
)

type Analyzer struct {
	Executor *Executor
}

func CreateAnalyzer() *Analyzer {
	analyzer := new(Analyzer)
	analyzer.Executor = ExecutorFactory()
	return analyzer
}

func (a *Analyzer) Analyze(pathToBinary string, binaryType int) binary.Binary {
	var file binary.Binary
	if isPortableExecutable(binaryType) {
		file = a.ProcessWindowsBinary(pathToBinary)
	} else if binaryType == TYPE_JAR {
		file = a.Jar(pathToBinary)
	}

	return file
}

func (a *Analyzer) CreateTemplateDirectory() {
	templateDir := a.Executor.TemplateDirectory
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		err := os.MkdirAll(templateDir, os.ModePerm)
		if err != nil {
			panic("Error creating template directory")
		}
	}
}

func (a *Analyzer) Jar(pathToJar string) *binary.JarBinary {
	jargoResult := Jargo(pathToJar)
	manifest := *jargoResult.Manifest

	jar := new(binary.JarBinary)
	jar.Architecture = ""
	jar.Filename = pathToJar
	jar.Dependencies = []binary.Dependency{}
	jar.Flags = []binary.Flag{}
	jar.Sections = []*pe.SectHeader{}
	jar.BuiltBy = manifest["Built-By"]
	jar.BuildJdk = manifest["Build-Jdk"]
	jar.CreatedBy = manifest["Created-By"]
	jar.ManifestVersion = manifest["Manifest-Version"]
	jar.MainClass = manifest["Main-Class"]
	fmt.Println(manifest["Class-Path"])
	jar.ClassPath = strings.Split(manifest["Class-Path"], " ")
	jar.JarAnalyzerTree = a.JarAnalyzer(pathToJar)

	return jar
}

func (a *Analyzer) ProcessWindowsBinary(pathToBinary string) *binary.PEBinary {
	bin := new(binary.PEBinary)
	objdump := a.ObjDump(pathToBinary, "a", "f", "x")
	pedumper := a.PEDumper(pathToBinary)
	peFile, err := pe.Open(pathToBinary)

	if err == nil {
		fileHeader, _ := peFile.FileHeader()
		bin.SectionNumber = fileHeader.NSection
		sectionHeaders, _ := peFile.SectHeaders()
		bin.Sections = sectionHeaders
	}

	bin.Filename = pathToBinary
	bin.Signature = getSignature(pedumper)
	bin.Size = getSize(pedumper)
	bin.Timestamp = getTimestamp(pedumper)
	bin.Time = getTime(bin.Timestamp)
	bin.Dependencies = getDllDependencies(objdump)
	bin.Architecture = getArchitecture(objdump)
	bin.ImportedFunctions = getImportedFunctions(pedumper)
	//bin.Sections = getSections(objdump)
	bin.Flags = getFlags(objdump)

	return bin
}

func getSize(dump *ObjDump) int64 {
	return getInteger(dump, "File size:\\s+(\\d+?) bytes")
}

func getTimestamp(dump *ObjDump) int64 {
	return getInteger(dump, "Timestamp:\\s+(\\d+)")
}

func getDllDependencies(dump *ObjDump) []binary.Dependency {
	depMatches := dump.FindAll("DLL Name: (.+?\\.dll)")
	dependencies := make([]binary.Dependency, len(depMatches))

	for index, element := range depMatches {
		dependencies[index] = binary.Dependency{Name: element[1]}
	}

	return dependencies
}

func getArchitecture(dump *ObjDump) string {
	return dump.FindAll("architecture: (.+?),")[0][1]
}

func getFlags(dump *ObjDump) []binary.Flag {
	flagsMatch := dump.FindAll("flags 0x[0-9a-f]+?:\\s+(([A-Z0-9_]+, )*[A-Z0-9_]+)\\s")
	flagStrings := strings.Split(flagsMatch[0][1], ", ")
	flags := make([]binary.Flag, len(flagStrings))

	for index, element := range flagStrings {
		flags[index] = binary.Flag{Name: element}
	}

	return flags
}

func getImportedFunctions(dump *ObjDump) []binary.Function {
	functionsMatch := dump.FindAll("Function Name \\(Hint\\):\\s+([^\\s]+)")
	functions := make([]binary.Function, len(functionsMatch))

	for index, element := range functionsMatch {
		functions[index] = binary.Function{Name: element[1]}
	}

	return functions
}

func getSignature(dump *ObjDump) string {
	signatureMatch := dump.Find("Signature:\\s+([^\\n]+)")
	return signatureMatch[1]
}

func getTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func getInteger(dump *ObjDump, regex string) int64 {
	timestampMatch := dump.Find(regex)
	timestamp, err := strconv.Atoi(timestampMatch[1])
	if err != nil {
		log.Fatal("Error convert timestamp to int")
		return -1
	}

	return int64(timestamp)
}

func isPortableExecutable(binaryType int) bool {
	return binaryType >= TYPE_EXE && binaryType <= TYPE_EFI
}
