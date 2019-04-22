package util

import (
	"fmt"
	"github.com/beevik/etree"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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

func IntToString(intValue int) string {
	return strconv.Itoa(intValue)
}

func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func GetLanguageByCompiler(compilerName string) string {
	lowerCase := strings.ToLower(compilerName)
	if strings.Contains(lowerCase, "gcc") {
		return "C/C++"
	}

	if strings.Contains(lowerCase, "jdk") || strings.Contains(lowerCase, "javac") {
		return "Java"
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

func LogIfError(err error, message string) {
	if err != nil {
		log.Println(message)
	}
}

