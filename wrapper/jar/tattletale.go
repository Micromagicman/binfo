package jar

import (
	"binfo/executable"
	osUtils "binfo/os"
	"binfo/wrapper"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Tattletale struct {
	wrapper.OnlyRun
	Tree *goquery.Document
}

func (tt *Tattletale) GetName() string {
	return "tattletale"
}

func (tt *Tattletale) LoadFile(pathToExecutable string) bool {
	if !tt.WasExecuted() {
		command := osUtils.Exec.TattletaleCommand(pathToExecutable)
		_, err := osUtils.Exec.Execute(command)
		if err != nil {
			return false
		}
		tt.MarkAsExecuted()
	}

	jarHtmlReport := osUtils.Exec.TemplateDirectory + "jar" + osUtils.Exec.Sep + filepath.Base(pathToExecutable) + ".html"
	htmlReport, err := os.Open(jarHtmlReport)
	if err != nil {
		return false
	}
	defer htmlReport.Close()

	document, err := goquery.NewDocumentFromReader(htmlReport)
	if err != nil {
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
	if manifestSelection == nil {
		return manifest
	}

	manifestNode := manifestSelection.Next()
	manifestHtml, err := manifestNode.Html()
	if err != nil {
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
	if requiresSelection == nil {
		return []string{}
	}

	requiresNode := requiresSelection.Next()
	requiresHtml, err := requiresNode.Html()

	if err != nil {
		return []string{}
	}
	return strings.Split(requiresHtml, "<br/>")
}

func (tt *Tattletale) getExports() []string {
	providesSelection := tt.findNodeWithText("td", "Provides")
	if providesSelection == nil {
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
	if nodes == nil {
		return nil
	}

	matches := nodes.FilterFunction(func(_ int, s *goquery.Selection) bool {
		return text == s.Text()
	})
	if matches.Size() == 0 {
		return nil
	}

	return matches.Eq(0)
}
