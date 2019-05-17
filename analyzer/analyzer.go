package analyzer

import (
	"binfo/executable"
	"binfo/util"
	"binfo/wrapper"
	"binfo/xml"
	"fmt"
	"github.com/H5eye/go-pefile"
	"github.com/decomp/exp/bin/elf"
	"github.com/decomp/exp/bin/pe"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	TYPE_UNKNOWN = -1
	TYPE_EXE     = 0

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
	TYPE_SO  = 10
	TYPE_AXF = 11
	TYPE_BIN = 12
	TYPE_ELF = 13
	TYPE_O   = 14
	TYPE_A   = 15
	TYPE_PRX = 16

	TYPE_JAR = 20
)

type Cache struct {
	Tattletale  bool
	JarAnalyzer bool
}

type Analyzer struct {
	Executor *Executor
	CompilerDetector *CompilerDetector
	Cache    *Cache
}

func CreateAnalyzer() *Analyzer {
	analyzer := new(Analyzer)
	analyzer.Cache = new(Cache)
	analyzer.Executor = ExecutorFactory()
	analyzer.CompilerDetector = CreateDetector("backend" + analyzer.Executor.Sep + "compiler_signatures.txt")
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
	jar.Filename, _ = filepath.Abs(pathToJar)
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

func (a *Analyzer) ProcessWindowsBinary(pathToExecutable string) (*executable.PortableExecutable, error) {
	fileStat, fileStatError := fileStat(pathToExecutable)
	if fileStatError != nil {
		return nil, fileStatError
	}

	binFile := new(executable.PortableExecutable)
	objDump := a.ObjDump(pathToExecutable, "x")

	binFile.Size = fileStat.Size()
	binFile.Filename, _ = filepath.Abs(pathToExecutable)
	binFile.Architecture = objDump.GetArchitecture()
	binFile.Exports = objDump.GetExports()
	binFile.Flags = objDump.GetFlags()
	binFile.Compiler = a.CompilerDetector.Detect(pathToExecutable)
	binFile.ProgrammingLanguage = util.GetLanguageByCompiler(binFile.Compiler)

	debugPe, err := wrapper.CreateDebugPeWrapper(pathToExecutable)
	if err != nil {
		log.Println("Cannot analyze " + pathToExecutable + " via decomp/exp/bin/pe library")
	} else {
		debugPe.Process(binFile)
	}

	peFile, err := pe.ParseFile(pathToExecutable)
	if err != nil {
		log.Println("Cannot analyze " + pathToExecutable + " via decomp/exp/bin/pe library")
	} else {
		binFile.Architecture = peFile.Arch.String()
	}

	memrevPE, err := wrapper.CreateMemrevPEWrapper(pathToExecutable)
	if err != nil {
		log.Println("Cannot analyze " + pathToExecutable + " via memrew/pe library")
	} else {
		memrevPE.Process(binFile)
	}

	peLoad, peError := pefile.Load(pathToExecutable)
	if peError != nil {
		log.Println("Cannot analyze " + pathToExecutable + " via H5eye/go-pefile library")
	} else {
		binFile.Timestamp = int64(peLoad.FileHeader.TimeDateStamp)
		binFile.Time = util.TimestampToTime(binFile.Timestamp)
	}

	return binFile, nil
}

func (a *Analyzer) ProcessLinuxBinary(pathToExecutable string) (*executable.ExecutableLinkable, error) {
	fileStat, fileStatError := fileStat(pathToExecutable)
	if fileStatError != nil {
		return nil, fileStatError
	}

	bin := new(executable.ExecutableLinkable)
	elfDump := a.ELFReader(pathToExecutable)

	bin.Filename, _ = filepath.Abs(pathToExecutable)
	bin.OperatingSystem = elfDump.GetOperatingSystem()
	bin.Size = fileStat.Size()
	bin.Format = elfDump.GetFormat()
	bin.Version = elfDump.GetVersion()
	bin.Endianess = elfDump.GetEndianess()
	bin.Compiler = a.CDetect(pathToExecutable)
	bin.ProgrammingLanguage = util.GetLanguageByCompiler(bin.Compiler)

	elfFile, err := elf.ParseFile(pathToExecutable)
	if err != nil {
		log.Println("Cannot analyze " + pathToExecutable + " via decomp/exp/bin/elf library")
	} else {
		bin.Architecture = elfFile.Arch.String()
		//bin.Imports = elfFile.Imports
		bin.Exports = elfFile.Exports
	}

	elfInfo, err := wrapper.CreateELFReaderWrapper(pathToExecutable)
	if err != nil {
		log.Println("Cannot analyze " + pathToExecutable + " via ELFReader")
	} else {
		elfInfo.Process(bin)
	}

	debugElf, err := wrapper.CreateDebugElfWrapper(pathToExecutable)
	if err != nil {
		log.Println("Cannot analyze " + pathToExecutable + " via ELFReader")
	} else {
		debugElf.Process(bin)
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
		util.LogIfError(err, "Error clear output directory")
	}
}

func (a *Analyzer) DeleteTemplateDirectory() {
	err := util.RemoveDirectory(a.Executor.TemplateDirectory)
	util.LogIfError(err, "Error removing template directory")
}

func (a *Analyzer) SaveResult(bin executable.Executable, outputDirectory string, path string) {
	xml.BuildXml(bin, outputDirectory+a.Executor.Sep+filepath.Base(path)+".xml")
}
