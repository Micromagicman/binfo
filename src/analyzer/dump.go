package analyzer

import "regexp"

type Dump interface {
	GetContent() string
}

type ObjDump struct {
	Content string
}

func (od *ObjDump) GetContent() string {
	return od.Content
}

func (od *ObjDump) Find(regExp string) []string {
	regex, _ := regexp.Compile(regExp)
	match := regex.FindStringSubmatch(od.Content)
	return match
}

func (od *ObjDump) FindAll(regExp string) [][]string {
	regex, _ := regexp.Compile(regExp)
	matches := regex.FindAllStringSubmatch(od.Content, -1)
	return matches
}