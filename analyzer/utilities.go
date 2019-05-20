package analyzer

import "binfo/wrapper"

type UtilitiesContainer struct {
	PE []wrapper.LibraryWrapper
	ELF []wrapper.LibraryWrapper
	JAR []wrapper.LibraryWrapper
}

func BuildUtilitiesContainer() *UtilitiesContainer {
	container := new(UtilitiesContainer)
	// Библиотеки для PE
	container.PE = []wrapper.LibraryWrapper{
		new(wrapper.DebugPe),
		new(wrapper.MemrevPE),
		new(wrapper.ObjDump),
	}
	// Библиотеки для ELF
	container.ELF = []wrapper.LibraryWrapper{
		new(wrapper.DebugElf),
		new(wrapper.ELFReader),
		new(wrapper.ELFReaderUtil),
		new(wrapper.CDetect),
	}
	// Библиотеки для Jar
	container.JAR = []wrapper.LibraryWrapper{

	}
	return container
}


//
//func (a *Analyzer) Tattletale(jarFilePath string) *Tattletale {
//	if !a.Cache.Tattletale {
//		command := a.Executor.TattletaleCommand(jarFilePath)
//		_, _ = a.Executor.Execute(command)
//		a.Cache.Tattletale = true
//	}
//	jarHtmlReport := "jar" + a.Executor.Sep + filepath.Base(jarFilePath) + ".html"
//	wrapper, err := CreateTattletaleWrapper(a.Executor.TemplateDirectory + jarHtmlReport)
//
//	if err != nil {
//		log.Println("Cannot analyze " + jarFilePath + " via Tattletale")
//		return nil
//	}
//
//	return wrapper
//}
//func (a *Analyzer) JarAnalyzer(pathToJar string) (*etree.Element, error) {
//	if !a.Cache.JarAnalyzer {
//		jarAnalyzerPath := a.Executor.AnalyzersPath + "jaranalyzer\\"
//		dir := filepath.Dir(pathToJar)
//		_, executeError := a.Executor.Execute(jarAnalyzerPath + "runxmlsummary.bat " + dir + " " + a.Executor.TemplateDirectory + "temp.xml")
//
//		if executeError != nil {
//			fmt.Println(executeError.Error())
//			return nil, executeError
//		}
//
//		a.Cache.JarAnalyzer = true
//	}
//
//	return getJarFileElement(a.Executor.TemplateDirectory+"temp.xml", pathToJar)
//}
//
//func getJarFileElement(pathToJarAnalyzerXml string, pathToJar string) (*etree.Element, error) {
//	doc := etree.NewDocument()
//	if err := doc.ReadFromFile(pathToJarAnalyzerXml); err != nil {
//		return nil, err
//	}
//
//	for _, jar := range doc.FindElements("//Jar") {
//		if strings.HasSuffix(pathToJar, jar.SelectAttr("name").Value) {
//			return jar.ChildElements()[0], nil // Summary
//		}
//	}
//
//	return nil, nil
//}
