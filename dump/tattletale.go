package dump

import (
	"binfo/executable"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Tattletale struct {
	Tree *goquery.Document
}

func CreateTattletaleWrapper(pathToHtml string) (*Tattletale, error) {
	fmt.Println(pathToHtml)
	tattletaleWrapper := new(Tattletale)
	htmlReport, err := os.Open(pathToHtml)
	if err != nil {
		return nil, err
	}
	defer htmlReport.Close()

	document, err := goquery.NewDocumentFromReader(htmlReport)
	if err != nil {
		return nil, err
	}

	tattletaleWrapper.Tree = document
	return tattletaleWrapper, nil
}

func (t *Tattletale) GetManifest() executable.JarManifest {
	manifestSelection := t.findNodeWithText("td", "Manifest")
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

func (t *Tattletale) GetRequires() []string {
	requiresSelection := t.findNodeWithText("td", "Requires")
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

func (t *Tattletale) GetProvides() []string {
	providesSelection := t.findNodeWithText("td", "Provides")
	if providesSelection == nil {
		return []string{}
	}

	return providesSelection.Next().
		Find("tr > td:first-child").
		Map(func(_ int, s *goquery.Selection) string {
			return s.Text()
		})
}

func (t *Tattletale) findNodeWithText(nodeName string, text string) *goquery.Selection {
	nodes := t.Tree.Find(nodeName)
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
