package analyzer

import "github.com/micromagicman/binary-info/wrapper"

type BinaryType int

const (
	TYPE_UNKNOWN BinaryType = iota
	TYPE_EXE
	// PE
	TYPE_DLL
	TYPE_OCX
	TYPE_SYS
	TYPE_SCR
	TYPE_DRV
	TYPE_CPL
	TYPE_EFI
	TYPE_LIB
	// ELF
	TYPE_SO
	TYPE_AXF
	TYPE_BIN
	TYPE_ELF
	TYPE_O
	TYPE_A
	TYPE_PRX
	// JAR
	TYPE_JAR
)

var extensions = map[string]BinaryType{
	"dll": TYPE_DLL,
	"ocx": TYPE_OCX,
	"sys": TYPE_SYS,
	"scr": TYPE_SCR,
	"drv": TYPE_DRV,
	"cpl": TYPE_CPL,
	"efi": TYPE_EFI,
	"lib": TYPE_LIB,
	"so":  TYPE_SO,
	"axf": TYPE_AXF,
	"bin": TYPE_BIN,
	"elf": TYPE_ELF,
	"o":   TYPE_O,
	"a":   TYPE_A,
	"prx": TYPE_PRX,
	"jar": TYPE_JAR,
}

type Analyzer struct {
	CompilerDetector *PECompilerDetector
	Utilities        *UtilitiesContainer
}

type PECompiler struct {
	Signature []interface{}
	EpOnly    bool
}

type PECompilerDetector struct {
	Compilers map[string]PECompiler
}

type UtilitiesContainer struct {
	Common []wrapper.LibraryWrapper
	PE     []wrapper.LibraryWrapper
	ELF    []wrapper.LibraryWrapper
	JAR    []wrapper.LibraryWrapper
}
