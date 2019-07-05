package jar

import (
	"github.com/micromagicman/binary-info/executable"
	osUtils "github.com/micromagicman/binary-info/os"
	"github.com/micromagicman/binary-info/wrapper"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Tattletale struct {
	wrapper.OnlyRun
	Tree *goquery.Document
}

func (tt *Tattletale) GetWindowsCommand() string {
	return "java"
}

func (tt *Tattletale) GetLinuxCommand() string {
	return "java"
}

func (tt *Tattletale) GetName() string {
	return "tattletale"
}

func (tt *Tattletale) LoadFile(pathToExecutable string) bool {
	arguments := []string{
		"-jar",
		osUtils.BackendDir + osUtils.Sep + "tattletale1.2.0.jar",
		filepath.Dir(pathToExecutable),
		osUtils.TemplateDir,
	}
	if !tt.WasExecuted() {
		if _, err := osUtils.Execute(tt, arguments...); nil != err {
			return false
		}
		tt.MarkAsExecuted()
	}
	jarHtmlReport := osUtils.TemplateDir + osUtils.Sep + "jar" + osUtils.Sep + filepath.Base(pathToExecutable) + ".html"
	htmlReport, err := os.Open(jarHtmlReport)
	if nil != err {
		return false
	}
	defer htmlReport.Close()
	document, err := goquery.NewDocumentFromReader(htmlReport)
	if nil != err {
		return false
	}
	tt.Tree = document
	return true
}

func (tt *Tattletale) Process(e executable.Executable) {
	jarFile := e.(*executable.JarExecutable)
	jarFile.Manifest = tt.getManifest()
	jarFile.Requires = tt.getImports()
	jarFile.Provides = tt.getExports()
}

func (tt *Tattletale) getManifest() executable.JarManifest {
	manifestSelection := tt.findNodeWithText("td", "Manifest")
	manifest := executable.JarManifest{}
	if nil == manifestSelection {
		return manifest
	}
	manifestNode := manifestSelection.Next()
	manifestHtml, err := manifestNode.Html()
	if nil != err {
		return manifest
	}
	manifestParameters := strings.Split(manifestHtml, "<br/>")
	var buffer, prevKey string
	for _, mp := range manifestParameters {
		keyValuePair := strings.Split(mp, ": ")
		if len(keyValuePair) < 2 {
			buffer += keyValuePair[0]
		} else {
			if "" != prevKey {
				manifest[prevKey] += buffer
				buffer = ""
			}
			prevKey = keyValuePair[0]
			manifest[keyValuePair[0]] = keyValuePair[1]
		}
		if "" != prevKey {
			manifest[prevKey] += buffer
			buffer = ""
		}
	}
	return manifest
}

func (tt *Tattletale) getImports() []string {
	requiresSelection := tt.findNodeWithText("td", "Requires")
	if nil == requiresSelection {
		return []string{}
	}
	requiresNode := requiresSelection.Next()
	requiresHtml, err := requiresNode.Html()
	if nil != err {
		return []string{}
	}
	return strings.Split(requiresHtml, "<br/>")
}

func (tt *Tattletale) getExports() []string {
	providesSelection := tt.findNodeWithText("td", "Provides")
	if nil == providesSelection {
		return []string{}
	}
	return providesSelection.Next().
		Find("tr > td:first-child").
		Map(func(_ int, s *goquery.Selection) string {
			return s.Text()
		})
}

func (tt *Tattletale) findNodeWithText(nodeName string, text string) *goquery.Selection {
	nodes := tt.Tree.Find(nodeName)
	if nil == nodes {
		return nil
	}
	matches := nodes.FilterFunction(func(_ int, s *goquery.Selection) bool {
		return text == s.Text()
	})
	if 0 == matches.Size() {
		return nil
	}
	return matches.Eq(0)
}
