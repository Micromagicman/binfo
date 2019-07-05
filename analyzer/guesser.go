package analyzer

func hasELFSignature(byteDump []byte) bool {
	return len(byteDump) >= 4 && byteDump[0] == 0x7F && byteDump[1] == 0x45 &&
		byteDump[2] == 0x4C && byteDump[3] == 0x46
}

func hasPESignature(byteDump []byte) bool {
	return len(byteDump) >= 2 && byteDump[0] == 0x4D && byteDump[1] == 0x5A
}

func hasJarSignature(byteDump []byte) bool {
	return len(byteDump) >= 4 && byteDump[0] == 0x50 && byteDump[1] == 0x4B &&
		(byteDump[2] == 0x03 && byteDump[3] == 0x04 || // Обычный архив
			byteDump[2] == 0x05 && byteDump[3] == 0x06 || // Пустой архив
			byteDump[2] == 0x07 && byteDump[3] == 0x08) // Составной архив
}
