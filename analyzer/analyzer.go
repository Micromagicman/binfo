package analyzer

import (
	"fmt"
	"github.com/micromagicman/binary-info/executable"
	osUtils "github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/util"
	"github.com/micromagicman/binary-info/wrapper"
	"io/ioutil"
	"log"
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

type Analyzer struct {
	CompilerDetector *PECompilerDetector
	Utilities *UtilitiesContainer
}

func CreateAnalyzer() *Analyzer {
	analyzer := new(Analyzer)
	analyzer.Utilities = BuildUtilitiesContainer()
	analyzer.CompilerDetector = CreateDetector(osUtils.BackendDir + osUtils.Sep + "compiler_signatures.txt")
	return analyzer
}

func (a *Analyzer) TryToAnalyze(unknownBinary string) (executable.Executable, error) {
	fileDump, err := ioutil.ReadFile(unknownBinary)
	if err != nil {
		return nil, err
	}

	if hasELFSignature(fileDump) {
		fmt.Println("Detect ELF signature")
		return a.ProcessExecutableLinkable(unknownBinary)
	} else if hasPESignature(fileDump) {
		fmt.Println("Detect PE signature")
		return a.ProcessPortableExecutable(unknownBinary)
	} else if hasJarSignature(fileDump) {
		fmt.Println("Detect ProcessJar signature")
		return a.ProcessJar(unknownBinary)
	}

	return nil, nil
}

func (a *Analyzer) Analyze(pathToBinary string, binaryType int) (executable.Executable, error) {
	var file executable.Executable
	var err error

	if binaryType == TYPE_EXE {
		file, err = a.TryToAnalyze(pathToBinary)
	} else if binaryType == TYPE_JAR {
		file, err = a.ProcessJar(pathToBinary)
	} else if isPortableExecutable(binaryType) {
		file, err = a.ProcessPortableExecutable(pathToBinary)
	} else if isElf(binaryType) {
		file, err = a.ProcessExecutableLinkable(pathToBinary)
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (a *Analyzer) ProcessJar(pathToJar string) (*executable.JarExecutable, error) {
	jar := new(executable.JarExecutable)
	jar.Compiler = "Javac"
	jar.ProgrammingLanguage = "Java"
	applyUtilities(pathToJar, jar, a.Utilities.Common, a.Utilities.JAR)
	return jar, nil
}

func (a *Analyzer) ProcessPortableExecutable(pathToExecutable string) (*executable.PortableExecutable, error) {
	binFile := new(executable.PortableExecutable)
	binFile.Compiler = a.CompilerDetector.Detect(pathToExecutable)
	binFile.ProgrammingLanguage = util.GetLanguageByCompiler(binFile.Compiler)
	applyUtilities(pathToExecutable, binFile, a.Utilities.Common, a.Utilities.PE)
	return binFile, nil
}

func (a *Analyzer) ProcessExecutableLinkable(pathToExecutable string) (*executable.ExecutableLinkable, error) {
	bin := new(executable.ExecutableLinkable)
	applyUtilities(pathToExecutable, bin, a.Utilities.Common, a.Utilities.ELF)
	return bin, nil
}

func applyUtilities(pathToExecutable string, e executable.Executable, utilities ...[]wrapper.LibraryWrapper) {
	for _, utilitiesPack := range utilities {
		for _, utility := range utilitiesPack {
			if utility.LoadFile(pathToExecutable) {
				utility.Process(e)
			} else {
				log.Println("Cannot analyze " + pathToExecutable + " via " + utility.GetName())
			}
		}
	}
}

func isPortableExecutable(binaryType int) bool {
	return binaryType >= TYPE_DLL && binaryType <= TYPE_LIB
}

func isElf(binaryType int) bool {
	return binaryType >= TYPE_SO && binaryType <= TYPE_PRX
}


