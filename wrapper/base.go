package wrapper

import (
	"github.com/micromagicman/binary-info/executable"
	"regexp"
)

type LibraryWrapper interface {
	GetName() string
	LoadFile(filename string) bool
	Process(e executable.Executable)
}

type Dump interface {
	GetContent() string
	Find(regExp string) []string
	FindAll(regExp string) [][]string
}

type BaseDump struct {
	Content string
}

type OnlyRun struct {
	Executed bool
}

func (or *OnlyRun) MarkAsExecuted() {
	or.Executed = true
}

func (or *OnlyRun) WasExecuted() bool {
	return or.Executed
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
	regex := regexp.MustCompile(regExp)
	matches := regex.FindAllStringSubmatch(bd.Content, -1)
	return matches
}

func (bd *BaseDump) SubDump(start int, end int) *BaseDump {
	subContent := ""
	if start > 0 && end < len(bd.Content) && start <= end {
		subContent = bd.Content[start:end]
	}
	return &BaseDump{subContent}
}

func Group(matches []string, requestIndex int) string {
	if len(matches) <= requestIndex {
		return ""
	}
	return matches[requestIndex]
}
