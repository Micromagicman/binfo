package os

import (
	"binfo/executable"
	"binfo/util"
	"binfo/xml"
	"os"
	"path/filepath"
)

func CreateTemplateDirectory() {
	templateDir := Exec.TemplateDirectory
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		err := util.CreateDirectory(templateDir)
		util.LogIfError(err, "Error creating template directory")
	} else {
		err := util.ClearDirectory(templateDir)
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
	err := util.RemoveDirectory(Exec.TemplateDirectory)
	util.LogIfError(err, "Error removing template directory")
}

func SaveResult(bin executable.Executable, outputDirectory string, path string) {
	xml.BuildXml(bin, outputDirectory+Exec.Sep+filepath.Base(path)+".xml")
}
