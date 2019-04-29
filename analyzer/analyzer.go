package analyzer

import (
	"binfo/executable"
	"binfo/util"
	"binfo/wrapper"
	"binfo/xml"
	"fmt"
	"github.com/mewrev/pe"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	TYPE_UNKNOWN = -1
	TYPE_EXE = 0

	// PE
	TYPE_DLL = 1
	TYPE_OCX = 2
	TYPE_SYS = 3
	TYPE_SCR = 4
	TYPE_DRV = 5
	TYPE_CPL = 6
	TYPE_EFI = 7
	TYPE_LIB = 8

	// ELF
	TYPE_SO = 10
	TYPE_AXF = 11
	TYPE_BIN = 12
	TYPE_ELF = 13
	TYPE_O = 14
	TYPE_A = 15
	TYPE_PRX = 16

	TYPE_JAR = 20
)

type Cache struct {
	Tattletale bool
	JarAnalyzer bool
}

type Analyzer struct {
	Executor *Executor
	Cache *Cache
}

func CreateAnalyzer() *Analyzer {
	analyzer := new(Analyzer)
	analyzer.Cache = new(Cache)
	analyzer.Executor = ExecutorFactory()
	return analyzer
}

func (a *Analyzer) TryToAnalyze(unknownBinary string) (executable.Executable, error) {
	fileDump, err := ioutil.ReadFile(unknownBinary)
	if err != nil {
		return nil, err
	}

	if hasELFSignature(fileDump) {
		fmt.Println("Detect ELF signature")
		return a.ProcessLinuxBinary(unknownBinary)
	} else if hasPESignature(fileDump) {
		fmt.Println("Detect PE signature")
		return a.ProcessWindowsBinary(unknownBinary)
	} else if hasJarSignature(fileDump) {
		fmt.Println("Detect Jar signature")
		return a.Jar(unknownBinary)
	}

	return nil, nil
}

func (a *Analyzer) Analyze(pathToBinary string, binaryType int) (executable.Executable, error) {
	var file executable.Executable
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

func (a *Analyzer) Jar(pathToJar string) (*executable.JarExecutable, error) {
	fileStat, fileStatError := fileStat(pathToJar)
	if fileStatError != nil {
		return nil, fileStatError
	}

	jar := new(executable.JarExecutable)
	jar.Size = fileStat.Size()
	jar.Filename = pathToJar
	jar.ProgrammingLanguage = "Java"

	tattletale := a.Tattletale(pathToJar)
	if tattletale != nil {
		jar.Manifest = tattletale.GetManifest()
		jar.Requires = tattletale.GetRequires()
		jar.Provides = tattletale.GetProvides()
	}

	jarAnalyzerTree, jarAnalyzerError := a.JarAnalyzer(pathToJar)
	if jarAnalyzerError != nil {
		log.Println("Cannot analyze file " + pathToJar + " via JarAnalyzer")
	} else {
		jar.JarAnalyzerTree = jarAnalyzerTree
	}

	return jar, nil
}

func (a *Analyzer) ProcessWindowsBinary(pathToBinary string) (*executable.PortableExecutable, error) {
	bin := new(executable.PortableExecutable)
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
	bin.Timestamp = peDumper.GetTimestamp()
	bin.Time = util.TimestampToTime(bin.Timestamp)
	bin.Dependencies = objDump.GetDependencies()
	bin.Architecture = objDump.GetArchitecture()
	bin.ImportedFunctions = peDumper.GetImportedFunctions()
	bin.ExportedFunctions = peDumper.GetExportedFunctions()
	bin.Flags = objDump.GetFlags()

	bin.Addresses = map[string]string{}
	bin.Addresses["EntryPoint"] = peDumper.GetEntryPointAddress()
	bin.Addresses["CodeSection"] = peDumper.GetCodeSectionAddress()
	bin.Addresses["DataSection"] = peDumper.GetDataSectionAddress()

	return bin, nil
}

func (a *Analyzer) ProcessLinuxBinary(pathToBinary string) (*executable.ExecutableLinkable, error) {
	fileStat, fileStatError := fileStat(pathToBinary)
	if fileStatError != nil {
		return nil, fileStatError
	}

	bin := new(executable.ExecutableLinkable)
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

func (a *Analyzer) InitOutputDirectory(outDir string) {
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

func (a *Analyzer) SaveResult(bin executable.Executable, outputDirectory string, path string) {
	xml.BuildXml(bin, outputDirectory + a.Executor.Sep + filepath.Base(path) + ".xml")
}