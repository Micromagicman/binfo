package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
)

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
	err := os.MkdirAll(name, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
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

func GetDirectoryFilePaths(directoryPath string) []string {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatal("Cannot read directory with binaries")
	}

	filePaths := make([]string, len(files))
	for i, file := range files {
		filePaths[i] = directoryPath + string(os.PathSeparator) + file.Name()
	}

	return filePaths
}

func LogIfError(err error, message string) {
	if err != nil {
		log.Println(message + ": " + err.Error())
	}
}
