Реализация метода "Извлечение информации из исполняемого файла программы"

На вход программе подается путь до директории с исполняемыми файлами. В результате работы в директории, указанной в качестве выходной,
для кажого исполняемого файла программа создает .xml файл в формате <имя_исполняемого_файла>.<расширение_исполняемого_файла>.xml с информацией, которую удалось извлечь.

## Запуск

### Cli
```go
go run cli/main.go -b=<directory with executables> -o=<output directory>
```

### Module
```go
import "github.com/micromagicman/binary-info/binfo"

func main() {
	binaries := "/some/binaries/dir"
	out := "/some/out/dir"
	binfo.ProcessFiles(binaries, out)
}
```

### Флаги

- **-d** - директория с исполняемыми файлами
- **-o** - директория с выходными xml-файлами