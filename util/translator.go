package util

import "strings"

var xmlMap = map[string]string{
	"executable":           "Исполняемый файл",
	"filename":             "Путь до файла",
	"magic":                "Сигнатура",
	"architecture":         "Архитектура",
	"programminglanguage":  "Язык программирования",
	"compiler":             "Компилятор",
	"compilationdate":      "Дата компиляции",
	"date":                 "Дата",
	"timestamp":            "UNIX Timestamp",
	"dependencies":         "Зависимости",
	"dependency":           "Зависимость",
	"size":                 "Размер файла",
	"manifest":             "Файл манифеста",
	"statistics":           "Статистика",
	"classcount":           "Количество классов",
	"abstractclasscount":   "Количество абстрактных классов",
	"packagecount":         "Количество пакетов",
	"packages":             "Пакеты",
	"package":              "Пакет",
	"class":                "Класс",
	"format":               "Формат",
	"type":                 "Тип",
	"importedfunctions":    "Импортируемые функции",
	"exportedfunctions":    "Экспортируемые функции",
	"function":             "Функция",
	"sections":             "Разделы",
	"section":              "Раздел",
	"flags":                "Флаги",
	"flag":                 "Флаг",
	"archiverversion":      "Версия архиватора",
	"builtby":              "Собрал",
	"specificationvendor":  "Разработчик спецификации",
	"specificationversion": "Версия спецификаии",
	"name":                 "Имя",
	"manifestversion":      "Версия манифеста",
	"createdby":            "Создано",
	"buildjdk":             "Версия JDK",
	"mainclass":            "Главный класс",
	"requires":             "Требования",
	"provides":             "Предоставляет",
}

func Translate(tagName string) string {
	if translation, has := xmlMap[strings.ToLower(tagName)]; has {
		return translation
	}
	return "-"
}
