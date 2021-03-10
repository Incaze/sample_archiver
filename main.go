package main

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	bufferByte := getInputFile()
	bufferForArchive, dictSlice, dictSize := getBufferForArchive(bufferByte)
	writeArchivedFile(bufferForArchive, dictSlice, dictSize)
	bufferArchived := getArchivedFile()
	decodedBuffer := decodeBuffer(bufferArchived)
	writeDecodedFile(decodedBuffer)
}
