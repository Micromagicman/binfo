package analyzer

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/logger"
	osUtils "github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/util"
	"github.com/micromagicman/binary-info/wrapper"
	"io"
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
	logger.Info("Processing " + binaryPath)
	binaryType := detectBinaryType(binaryPath)
	if binaryType != TYPE_UNKNOWN {
		binFile, err = a.AnalyzeNormal(binaryPath, binaryType)
	} else {
		// попытка проанализировать файл без расширения
		logger.Warning("Unknown executable type for file " + binaryPath)
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
	fileDump, err := util.FileFirstBytes(unknownBinary, 4)
	if nil != err {
		return nil, err
	}
	if HasELFSignature(fileDump) {
		logger.Info("Detect ELF signature for " + unknownBinary)
		return a.ProcessExecutableLinkable(unknownBinary)
	}
	if HasPESignature(fileDump) {
		logger.Info("Detect PE signature for " + unknownBinary)
		return a.ProcessPortableExecutable(unknownBinary)
	}
	if HasJarSignature(fileDump) {
		logger.Info("Detect jar signature for " + unknownBinary)
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
	fmt.Println("Process " + pathToJar)
	//list := util.FindInnerJars(pathToJar)
	//for _, jarChild := range list {
	//	log.Println("Found jar file " + filepath.Base(jarChild) + " inside " + pathToJar)
	//	jarChildExecutable, err := a.ProcessJar(jarChild)
	//	if nil != err {
	//		log.Println("Cannot analyze file " + jarChild)
	//	} else {
	//		jar.Children = append(jar.Children, jarChildExecutable)
	//	}
	//}
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
				logger.Warning("Cannot analyze " + pathToExecutable + " via " + utility.GetName())
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

func FindAndCreateInnerJars(binariesDirectory string) []string {
	var innerJarPaths []string
	for _, file := range util.GetDirectoryFilePaths(binariesDirectory) {
		fileBytes, _ := util.FileFirstBytes(file, 4)
		if HasJarSignature(fileBytes) {
			innerJarPaths = append(innerJarPaths, innerJars(file)...)
		}
	}
	return innerJarPaths
}

func innerJars(jarFilePath string) []string {
	read, err := zip.OpenReader(jarFilePath)
	if nil != err {
		log.Fatal(err.Error())
	}
	var paths []string
	for _, f := range read.File {
		fileBytes, err := getJarInsideJarFileBytes(f)
		if nil != err {
			continue
		}
		fileName := filepath.Base(jarFilePath) + "__" + filepath.Base(f.Name)
		savePath := filepath.Join(filepath.Dir(jarFilePath), fileName)
		if err = ioutil.WriteFile(savePath, fileBytes, 0664); nil != err {
			log.Println(err.Error())
			continue
		}
		paths = append(paths, savePath)
	}
	for _, jp := range paths {
		paths = append(paths, innerJars(jp)...)
	}
	return paths
}

func getJarInsideJarFileBytes(f *zip.File) ([]byte, error) {
	file, err := f.Open()
	if nil != err {
		return nil, err
	}
	defer file.Close()
	buffer := new(bytes.Buffer)
	_, err = io.CopyN(buffer, file, int64(f.UncompressedSize64))
	if nil != err || !HasJarSignature(buffer.Bytes()[:4]) {
		return nil, errors.New("File " + f.Name + " is not jar")
	}
	return buffer.Bytes(), nil
}
