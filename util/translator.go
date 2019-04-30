package util

var xmlMap = map[string]string {
	"Executable": "Исполняемый файл",
	"Filename": "Путь до файла",
	"Magic": "Сигнатура",
	"Architecture": "Архитектура",
	"ProgrammingLanguage": "Язык программирования",
	"Compiler": "Компилятор",
	"CompilationDate": "Дата компиляции",
	"Date": "Дата",
	"Timestamp": "UNIX Timestamp",
	"Dependencies": "Зависимости",
	"Dependency": "Зависимость",
	"Size": "Размер файла",
	"Manifest": "Файл манифеста",
	"Statistics": "Статистика",
	"ClassCount": "Количество классов",
	"AbstractClassCount": "Количество абстрактных классов",
	"PackageCount": "Количество пакетов",
	"Packages": "Пакеты",
	"Package": "Пакет",
	"Class": "Класс",
}

func Translate(tagName string) string {
	if translation, has := xmlMap[tagName]; has {
		return translation
	}
	return "-"
}