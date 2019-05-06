Реализация метода "Извлечение информации из исполняемого кода программы"

## Запуск

```
go run main.go -b=<Директория с бинарниками> -o=<Директория с результатом работы>
```

## Перевод xml тегов
```go
import "binfo/util"

tagName := "Class"
translated := util.Translate(tagName)
```