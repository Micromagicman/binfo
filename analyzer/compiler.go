package analyzer

import (
	"encoding/hex"
	"github.com/micromagicman/binary-info/executable"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type PECompiler struct {
	Signature []interface{}
	EpOnly bool
}

type PECompilerDetector struct {
	Compilers map[string]PECompiler
}


func (cd *PECompilerDetector) Detect(pathToExecutable string) string {
	file, _ := os.Open(pathToExecutable)
	maxMatchCount := 0
	bestMatchCompilerName := executable.DEFAULT_VALUE
	for compilerName, compiler := range cd.Compilers {
		currentMatchCount := compiler.Match(file)
		if currentMatchCount > maxMatchCount {
			maxMatchCount = currentMatchCount
			bestMatchCompilerName = compilerName
		}
	}
	return bestMatchCompilerName
}

func (c *PECompiler) Match(file *os.File) int {
	fileBytes := make([]byte, len(c.Signature))
	file.Read(fileBytes)
	file.Seek(0, 0)
	matchCount := 0
	for i, fb := range fileBytes {
		if sb, ok := c.Signature[i].(byte); ok {
			if fb != sb {
				return 0
			}
			matchCount++
		}
	}
	return matchCount
}

func CreateDetector(pathToDatabase string) *PECompilerDetector {
	signatures := parseSignatures(pathToDatabase)
	detector := new(PECompilerDetector)
	detector.Compilers = map[string]PECompiler{}
	// [1] - Имя компилятора
	// [2] - сигнатура
	// [5] - EpOnly
	for _, s := range signatures {
		epOnly := false
		if "true" == s[5] {
			epOnly = true
		}
		compilerName := s[1]
		signature := createByteSignature(s[2])
		detector.Compilers[compilerName] = PECompiler{signature, epOnly}
	}
	return detector
}

func createByteSignature(stringSignature string) []interface{} {
	arrayOfBytes := strings.Split(stringSignature, " ")
	var bytes []interface{}
	for _, b := range arrayOfBytes {
		if "??" == b {
			bytes = append(bytes, "??")
		} else {
			decodedByte, _ := hex.DecodeString(b)
			bytes = append(bytes, decodedByte[0])
		}
	}
	return bytes
}

func parseSignatures(pathToDatabase string) [][]string {
	fileBytes, err := ioutil.ReadFile(pathToDatabase)
	if err != nil {
		return [][]string{}
	}
	fileContent := string(fileBytes)
	regex := regexp.MustCompile(`\[(.*?)\]\s+?signature\s*=\s*(.*?)(\s+\?\?)*\s*ep_only\s*=\s*(\w+)(?:\s*section_start_only\s*=\s*(\w+)|)`)
	return regex.FindAllStringSubmatch(fileContent, -1)
}


