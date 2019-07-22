package xml

import (
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/logger"
	"os"

	"github.com/beevik/etree"
)

func BuildXml(bin executable.Executable, outputFilename string) {
	file, _ := os.Create(outputFilename)
	doc := etree.NewDocument()
	bin.BuildXml(doc)
	doc.Indent(4)
	_, err := doc.WriteTo(file)
	if nil != err {
		logger.Error("Error creating xml output for file " + outputFilename)
	}
}
