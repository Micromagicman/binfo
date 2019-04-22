package analyzer

import (
	"binary"
	"github.com/mewrev/pe"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"util"
	"wrapper"
	"xml"
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
	TYPE_A = 15
	TYPE_PRX = 16
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

func (a *Analyzer) TryToAnalyze(unknownBinary string) (binary.Binary, error) {
	fileDump, err := ioutil.ReadFile(unknownBinary)
	if err != nil {
		return nil, err
	}

	if hasELFSignature(fileDump) {
		return a.ProcessLinuxBinary(unknownBinary)
	} else if hasPESignature(fileDump) {
		return a.ProcessWindowsBinary(unknownBinary)
	} else if hasJarSignature(fileDump) {
		return a.Jar(unknownBinary)
	}

	return nil, nil
}

func (a *Analyzer) Analyze(pathToBinary string, binaryType int) (binary.Binary, error) {
	var file binary.Binary
	var err error

	if binaryType == TYPE_EXE {
		file, err = a.TryToAnalyze(pathToBinary)
	} else if binaryType == TYPE_JAR {
		file, err = a.Jar(pathToBinary)
	} else if isPortableExecutable(binaryType) {
		file, err = a.ProcessWindowsBinary(pathToBinary)
	} else if isElf(binaryType) {
		file, err = a.ProcessLinuxBinary(pathToBinary)
	}

	if err != nil {
		return nil, err
	}

	return file, nil
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
	jar.ProgrammingLanguage = "Java"
	jar.Dependencies = []binary.Dependency{}
	jar.Flags = []binary.Flag{}
	jar.Sections = []*pe.SectHeader{}
	jar.BuiltBy = manifest["Built-By"]
	jar.BuildJdk = manifest["Build-Jdk"]
	jar.CreatedBy = manifest["Created-By"]
	jar.ManifestVersion = manifest["Manifest-Version"]
	jar.MainClass = manifest["Main-Class"]
	jar.ClassPath = strings.Split(manifest["Class-Path"], " ")
	jarAnalyzerTree, jarAnalyzerError := a.JarAnalyzer(pathToJar)

	if jarAnalyzerError != nil {
		log.Fatal("Cannot analyze file " + pathToJar + " via JarAnalyzer")
	} else {
		jar.JarAnalyzerTree = jarAnalyzerTree

	}

	return jar, nil
}

func (a *Analyzer) ProcessWindowsBinary(pathToBinary string) (*binary.PEBinary, error) {
	bin := new(binary.PEBinary)
	objDump := a.ObjDump(pathToBinary, "x")
	peDumper := a.PEDumper(pathToBinary)
	peFile, peError := pe.Open(pathToBinary)

	if peError == nil {
		fileHeader, _ := peFile.FileHeader()
		bin.SectionNumber = fileHeader.NSection
		sectionHeaders, _ := peFile.SectHeaders()
		bin.Sections = sectionHeaders
	}

	bin.Filename = pathToBinary
	bin.Size = peDumper.GetSize()
	bin.EntryPointAddress = peDumper.GetEntryPointAddress()
	bin.Timestamp = peDumper.GetTimestamp()
	bin.Time = util.TimestampToTime(bin.Timestamp)
	bin.Dependencies = objDump.GetDependencies()
	bin.Architecture = objDump.GetArchitecture()
	bin.ImportedFunctions = peDumper.GetImportedFunctions()
	bin.ExportedFunctions = peDumper.GetExportedFunctions()
	bin.Flags = objDump.GetFlags()

	return bin, nil
}

func (a *Analyzer) ProcessLinuxBinary(pathToBinary string) (*binary.ELFBinary, error) {
	fileStat, fileStatError := fileStat(pathToBinary)
	if fileStatError != nil {
		return nil, fileStatError
	}

	bin := new(binary.ELFBinary)
	elfDump := a.ELFReader(pathToBinary)

	bin.Filename = pathToBinary
	bin.OperatingSystem = elfDump.GetOperatingSystem()
	bin.Size = fileStat.Size()
	bin.Format = elfDump.GetFormat()
	bin.Version = elfDump.GetVersion()
	bin.Endianess = elfDump.GetEndianess()
	bin.Type = elfDump.GetType()
	compilerBytes, elfInfoError := a.Executor.Execute(a.Executor.ELFInfoCommand(pathToBinary))

	if elfInfoError == nil {
		bin.Compiler = string(compilerBytes)
		bin.ProgrammingLanguage = util.GetLanguageByCompiler(bin.Compiler)
	}

	elfInfo, err := wrapper.CreateELFReader(pathToBinary)
	if err == nil {
		bin.Sections = elfInfo.GetSections()
		bin.ImportedFunctions = elfInfo.GetImportedFunctions()
	}

	return bin, nil
}

func isPortableExecutable(binaryType int) bool {
	return binaryType >= TYPE_DLL && binaryType <= TYPE_LIB
}

func isElf(binaryType int) bool {
	return binaryType >= TYPE_SO && binaryType <= TYPE_PRX
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

func (a *Analyzer) CreateTemplateDirectory() {
	templateDir := a.Executor.TemplateDirectory
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		err := util.CreateDirectory(templateDir)
		util.LogIfError(err, "Error creating template directory")
	}
}

func (a *Analyzer) InitOutputDirectory() {
	outDir := "out"
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := util.CreateDirectory(outDir)
		util.LogIfError(err, "Error creating output directory")
	} else {
		err := util.ClearDirectory(outDir)
		util.LogIfError(err,"Error clear output directory")
	}
}

func (a *Analyzer) DeleteTemplateDirectory() {
	err := util.RemoveDirectory(a.Executor.TemplateDirectory)
	util.LogIfError(err, "Error removing template directory")
}

func (a *Analyzer) SaveResult(bin binary.Binary, path string) {
	xml.BuildXml(bin, "out" + a.Executor.Sep + filepath.Base(path) + ".xml")
}