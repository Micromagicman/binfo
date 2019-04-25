package xml

import (
	"binfo/binary"
	"github.com/beevik/etree"
	"log"
	"os"
)

func BuildXml(bin binary.Binary, outputFilename string) {
	file, _ := os.Create(outputFilename)
	doc := etree.NewDocument()
	bin.BuildXml(doc)

	doc.Indent(4)
	_, err := doc.WriteTo(file)
	if err != nil {
		log.Println("Error creating xml output for file " + outputFilename)
	}
}