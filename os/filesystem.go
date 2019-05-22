package os

import (
	"binfo/executable"
	"binfo/util"
	"binfo/xml"
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

func SaveResult(bin executable.Executable, outputDirectory string, path string) {
	xml.BuildXml(bin, outputDirectory+Sep+filepath.Base(path)+".xml")
}
