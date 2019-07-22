package os

import (
	"github.com/micromagicman/binary-info/executable"
	"github.com/micromagicman/binary-info/util"
	"github.com/micromagicman/binary-info/xml"
	"log"
	"os"
	"path/filepath"
)

func CreateTemplateDirectory() {
	if _, err := os.Stat(TemplateDir); os.IsNotExist(err) {
		err := util.CreateDirectory(TemplateDir)
		util.LogIfError(err, "Error creating template directory")
	} else {
		err := util.ClearDirectory(TemplateDir)
		util.LogIfError(err,"Error clear template directory")
	}
}

func InitOutputDirectory(outDir string) {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := util.CreateDirectory(outDir)
		util.LogIfError(err, "Error creating output directory")
	} else {
		err := util.ClearDirectory(outDir)
		util.LogIfError(err, "Error clear output directory")
	}
}

func DeleteTemplateDirectory() {
	err := util.RemoveDirectory(TemplateDir)
	util.LogIfError(err, "Error removing template directory")
}

func DeleteFiles(files... string) {
	for _, f := range files {
		if err := os.Remove(f); nil != err {
			log.Println(err.Error())
		}
	}
}

func SaveResult(bin executable.Executable, outputPath string) {
	_ = util.CreateDirectory(filepath.Dir(outputPath))
	xml.BuildXml(bin, outputPath)
}