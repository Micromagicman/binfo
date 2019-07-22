package util

import (
	"fmt"
	"github.com/micromagicman/binary-info/logger"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
)

var hashString = "abcdef1234567890"
var languageToCompiler = map[string][]string{
	"Golang":  {"go", "golang"},
	"C/C++":   {"gcc", "clang", "visual c++", "borland c++", "mingw"},
	"Java":    {"javac"},
	"Haskell": {"ghc"},
	"Rust":    {"rust"},
	"OCaml":   {"ocaml"},
	"Pascal":  {"fpc"},
	"C":       {"tcc"},
	"Delphi":  {"borland delphi"},
}

func BuildNodeWithText(nodeName string, nodeContent string) *etree.Element {
	node := etree.NewElement(nodeName)
	node.CreateText(nodeContent)
	return node
}

func GetOptionalStringValue(value string, negativeResult string) string {
	if value == "" {
		return negativeResult
	}
	return value
}

func UInt64ToString(intValue uint64) string {
	return strconv.FormatUint(intValue, 10)
}

func Int64ToString(intValue int64) string {
	return strconv.FormatInt(intValue, 10)
}

func Uint32ToString(intValue uint32) string {
	return fmt.Sprint(intValue)
}

func Int64ToHex(intValue int64) string {
	return "0x" + strconv.FormatInt(intValue, 16)
}

func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func GetLanguageByCompiler(compilerName string) string {
	lowerCase := strings.ToLower(compilerName)
	for lang, compilers := range languageToCompiler {
		for _, compiler := range compilers {
			if strings.Contains(lowerCase, compiler) {
				return lang
			}
		}
	}
	return "Unknown"
}

func CreateDirectory(name string) error {
	return os.MkdirAll(name, 0664)
}

func RemoveDirectory(name string) error {
	return os.RemoveAll(name)
}

func ClearDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func CreateOutputPath(executablePath string, outputPath string) string {
	return filepath.Join(outputPath, executablePath + ".xml")
}

func GetDirectoryFilePaths(directoryPath string) []string {
	var filePaths []string
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if nil != err {
			return err
		}
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	if nil != err {
		log.Fatal("Cannot read directory with binaries")
	}
	return filePaths
}

func FileFirstBytes(filePath string, countOfBytes int64) ([]byte, error) {
	file, err := os.Open(filePath)
	if nil != err {
		return nil, err
	}
	defer file.Close()
	buffer := make([]byte, countOfBytes)
	n, err := io.ReadFull(file, buffer)
	if nil != err || int64(n) != countOfBytes {
		return nil, err
	}
	return buffer, nil
}

func LogIfError(err error, message string) {
	if err != nil {
		logger.Error(message + ": " + err.Error())
	}
}

func StringInSlice(el string, list []string) bool {
	for _, item := range list {
		if item == el {
			return true
		}
	}
	return false
}

func AppendIfNotExists(el string, list[]string) []string {
	if !StringInSlice(el, list) {
		list = append(list, el)
	}
	return list
}

func MapKeys(m map[interface{}]interface{}) []interface{} {
	var keys []interface{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapValues(m map[interface{}]interface{}) []interface{} {
	var values []interface{}
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func RandomHash(size int) []byte {
	hash := make([]byte, size)
	for i := range hash {
		hash[i] = hashString[rand.Intn(len(hashString))]
	}
	return hash
}