package wrapper

import "binfo/executable"

type LibraryWrapper interface {
	Process(bin executable.Executable)
}
