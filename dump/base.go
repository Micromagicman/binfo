package dump

import (
	"regexp"
	"strconv"
)

type Dump interface {
	GetContent() string
	Find(regExp string) []string
	FindAll(regExp string) [][]string
}

type BaseDump struct {
	Content string
}

func (bd *BaseDump) GetContent() string {
	return bd.Content
}

func (bd *BaseDump) Find(regExp string) []string {
	regex, _ := regexp.Compile(regExp)
	match := regex.FindStringSubmatch(bd.Content)
	return match
}

func (bd *BaseDump) FindAll(regExp string) [][]string {
	regex, _ := regexp.Compile(regExp)
	matches := regex.FindAllStringSubmatch(bd.Content, -1)
	return matches
}

func GetInteger(dump Dump, regex string) int64 {
	timestampMatch := dump.Find(regex)
	timestamp, err := strconv.Atoi(Group(timestampMatch, 1))
	if err != nil {
		return -1
	}

	return int64(timestamp)
}

func Group(matches []string, requestIndex int) string {
	if len(matches) < requestIndex {
		return ""
	}
	return matches[requestIndex]
}
