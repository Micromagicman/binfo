package util

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
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

func Int64ToString(intValue int64) string {
	return strconv.FormatInt(intValue, 10)
}

func Uint32ToString(intValue uint32) string {
	return fmt.Sprint(intValue)
}

func IntToString(intValue int) string {
	return strconv.Itoa(intValue)
}