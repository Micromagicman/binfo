package xml

import (
	"binary"
	"github.com/beevik/etree"
	"log"
	"os"
)

func BuildXml(bin *binary.PEBinary, outputFilename string) {
	file, _ := os.Create(outputFilename)
	doc := etree.NewDocument()
	bin.ToXml(doc)

	doc.Indent(4)
	_, err := doc.WriteTo(file)
	if err != nil {
		log.Fatal("Error when create xml")
	}
}