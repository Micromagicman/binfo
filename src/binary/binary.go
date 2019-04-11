package binary

type Section struct {
	Idx int
	Name string
	Size int
}

type Dependency struct {
	Name string
}

type Flag struct {
	Name string
}

type BinaryFile struct {
	Filename string
	Architecture string
	Dependencies []Dependency
	Flags []Flag
	Sections []Section
}
