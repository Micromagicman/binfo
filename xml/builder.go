package xml

import (
	"binfo/executable"
	"log"
	"os"

	"github.com/beevik/etree"
)

func BuildXml(bin executable.Executable, outputFilename string) {
	file, _ := os.Create(outputFilename)
	doc := etree.NewDocument()
	bin.BuildXml(doc)

	doc.Indent(4)
	_, err := doc.WriteTo(file)
	if err != nil {
		log.Println("Error creating xml output for file " + outputFilename)
	}
}
