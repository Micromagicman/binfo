package analyzer

import (
	"binary"
	"github.com/mewrev/pe"
	"os"
	"strings"
	"util"
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
	TYPE_LIB = 8
	TYPE_SO = 10
	TYPE_AXF = 11
	TYPE_BIN = 12
	TYPE_ELF = 13
	TYPE_O = 14
	TYPE_PRX = 15
	TYPE_JAR = 20
)

type Analyzer struct {
	Executor *Executor
}

func CreateAnalyzer() *Analyzer {
	analyzer := new(Analyzer)
	analyzer.Executor = ExecutorFactory()
	return analyzer
}

func (a *Analyzer) Analyze(pathToBinary string, binaryType int) (binary.Binary, error) {
	var file binary.Binary
	var err error
	if isPortableExecutable(binaryType) {
		file, err = a.ProcessWindowsBinary(pathToBinary)
	} else if isElf(binaryType) {
		file, err = a.ProcessLinuxBinary(pathToBinary)
	} else if binaryType == TYPE_JAR {
		file, err = a.Jar(pathToBinary)
	}

	if err != nil {
		return nil, err
	}

	return file, nil
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

func (a *Analyzer) Jar(pathToJar string) (*binary.JarBinary, error) {
	fileStat, fileStatError := fileStat(pathToJar)
	if fileStatError != nil {
		return nil, fileStatError
	}

	jargoResult := Jargo(pathToJar)
	manifest := *jargoResult.Manifest

	jar := new(binary.JarBinary)
	jar.Architecture = ""
	jar.Size = fileStat.Size()
	jar.Filename = pathToJar
	jar.Dependencies = []binary.Dependency{}
	jar.Flags = []binary.Flag{}
	jar.Sections = []*pe.SectHeader{}
	jar.BuiltBy = manifest["Built-By"]
	jar.BuildJdk = manifest["Build-Jdk"]
	jar.CreatedBy = manifest["Created-By"]
	jar.ManifestVersion = manifest["Manifest-Version"]
	jar.MainClass = manifest["Main-Class"]
	jar.ClassPath = strings.Split(manifest["Class-Path"], " ")
	jar.JarAnalyzerTree = a.JarAnalyzer(pathToJar)

	return jar, nil
}

func (a *Analyzer) ProcessWindowsBinary(pathToBinary string) (*binary.PEBinary, error) {
	bin := new(binary.PEBinary)
	objDump := a.ObjDump(pathToBinary, "a", "f", "x")
	peDumper := a.PEDumper(pathToBinary)
	peFile, peError := pe.Open(pathToBinary)

	if peError == nil {
		fileHeader, _ := peFile.FileHeader()
		bin.SectionNumber = fileHeader.NSection
		sectionHeaders, _ := peFile.SectHeaders()
		bin.Sections = sectionHeaders
	}

	bin.Filename = pathToBinary
	bin.Signature = peDumper.GetSignature()
	bin.Size = peDumper.GetSize()
	bin.Timestamp = peDumper.GetTimestamp()
	bin.Time = util.TimestampToTime(bin.Timestamp)
	bin.Dependencies = objDump.GetDependencies()
	bin.Architecture = objDump.GetArchitecture()
	bin.ImportedFunctions = peDumper.GetImportedFunctions()
	bin.Flags = objDump.GetFlags()

	return bin, nil
}

func (a *Analyzer) ProcessLinuxBinary(pathToBinary string) (*binary.ELFBinary, error) {
	bin := new(binary.ELFBinary)
	elfDump := a.ELFReader(pathToBinary)

	bin.OperatingSystem = elfDump.GetOperatingSystem()
	bin.Format = elfDump.GetFormat()
	bin.Version = elfDump.GetVersion()
	bin.Endianess = elfDump.GetEndianess()
	bin.Type = elfDump.GetType()

	return bin, nil
}

func isPortableExecutable(binaryType int) bool {
	return binaryType >= TYPE_EXE && binaryType <= TYPE_LIB
}

func isElf(binaryType int) bool {
	return binaryType == TYPE_EXE || (binaryType >= TYPE_SO && binaryType <= TYPE_PRX)
}

func fileStat(path string) (os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return info, nil
}
