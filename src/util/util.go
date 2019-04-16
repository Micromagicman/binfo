package util

import "github.com/beevik/etree"

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
