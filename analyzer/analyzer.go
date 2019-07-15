package analyzer

import (
	"github.com/pkg/errors"
	"github.com/micromagicman/binary-info/executable"
	osUtils "github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/util"
	"github.com/micromagicman/binary-info/wrapper"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
)

func CreateAnalyzer() *Analyzer {
	return &Analyzer{
		CreateCompilerDetector(path.Join(osUtils.BackendDir, "compiler_signatures.txt")),
		BuildUtilitiesContainer(),
	}
}

func (a *Analyzer) Analyze(binaryPath string) (executable.Executable, error) {
	var binFile executable.Executable
	var err error
	log.Println("Processing " + binaryPath)
	binaryType := detectBinaryType(binaryPath)
	if binaryType != TYPE_UNKNOWN {
		binFile, err = a.AnalyzeNormal(binaryPath, binaryType)
	} else {
		// попытка проанализировать файл без расширения
		log.Println("Unknown executable type for file " + binaryPath)
		binFile, err = a.AnalyzeUnknown(binaryPath)
		if nil != err {
			return nil, err
		}
	}
	if nil != err {
		return nil, err
	} else {
		return binFile, nil
	}
}

func (a *Analyzer) AnalyzeUnknown(unknownBinary string) (executable.Executable, error) {
	fileDump, err := ioutil.ReadFile(unknownBinary)
	if nil != err {
		return nil, err
	}
	if hasELFSignature(fileDump) {
		log.Println("Detect ELF signature for " + unknownBinary)
		return a.ProcessExecutableLinkable(unknownBinary)
	}
	if hasPESignature(fileDump) {
		log.Println("Detect PE signature for " + unknownBinary)
		return a.ProcessPortableExecutable(unknownBinary)
	}
	if hasJarSignature(fileDump) {
		log.Println("Detect jar signature for " + unknownBinary)
		return a.ProcessJar(unknownBinary)
	}
	return nil, errors.New("Unknown executable type for " + unknownBinary)
}

func (a *Analyzer) AnalyzeNormal(pathToBinary string, binaryType BinaryType) (executable.Executable, error) {
	var file executable.Executable
	var err error
	if binaryType == TYPE_EXE {
		file, err = a.AnalyzeUnknown(pathToBinary)
	} else if binaryType == TYPE_JAR {
		file, err = a.ProcessJar(pathToBinary)
	} else if isPortableExecutable(binaryType) {
		file, err = a.ProcessPortableExecutable(pathToBinary)
	} else if isElf(binaryType) {
		file, err = a.ProcessExecutableLinkable(pathToBinary)
	}
	if nil != err {
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

func isPortableExecutable(binaryType BinaryType) bool {
	return binaryType >= TYPE_DLL && binaryType <= TYPE_LIB
}

func isElf(binaryType BinaryType) bool {
	return binaryType >= TYPE_SO && binaryType <= TYPE_PRX
}

func detectBinaryType(binaryPath string) BinaryType {
	extension := filepath.Ext(binaryPath)
	if "" != extension {
		extension = extension[1:]
	}
	if value, ok := extensions[extension]; ok {
		return value
	}
	return TYPE_UNKNOWN
}


