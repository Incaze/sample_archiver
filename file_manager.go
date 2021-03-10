package main

import (
	"io/ioutil"
	"os"
)

const (
	INPUT_FILE        = "dreamshed.txt"
	ARCHIVE_EXTENSION = "32a"
)

func writeArchivedFile(buffer []byte, dictSlice [][]byte, dictSize []byte) {
	outputFile := INPUT_FILE + "." + ARCHIVE_EXTENSION
	encodedExtension := []byte(ARCHIVE_EXTENSION)
	file, err := os.Create(outputFile)
	checkError(err)
	defer file.Close()
	file.Write(encodedExtension)
	file.Write(dictSize)
	for _, bytes := range dictSlice {
		file.Write(bytes)
	}
	file.Write(buffer)
}

func writeDecodedFile(buffer [][]byte) {
	outputFile := "decoded_" + INPUT_FILE
	file, err := os.Create(outputFile)
	checkError(err)
	defer file.Close()
	for _, bytes := range buffer {
		file.Write(bytes)
	}
}

func getArchivedFile() []byte {
	buffer, err := ioutil.ReadFile(INPUT_FILE + "." + ARCHIVE_EXTENSION)
	checkError(err)
	return buffer
}

func getInputFile() [][]byte {
	buffer, err := ioutil.ReadFile(INPUT_FILE)
	checkError(err)
	var result [][]byte
	bufferLen := len(buffer)
	for i := 0; i < bufferLen; {
		batch := make([]byte, 0, BATCH_SIZE)
		left := bufferLen - (i + 2)
		if left > 0 {
			batch = append(batch, buffer[i], buffer[i+1], buffer[i+2])
		} else {
			for j := i; j < bufferLen; j++ {
				batch = append(batch, buffer[j])
			}
			batchLen := len(batch)
			for j := batchLen; j < BATCH_SIZE; j++ {
				batch = append(batch, MANAGE_BYTE)
			}
		}
		if len(batch) == 0 {
			break
		}
		result = append(result, batch)
		i += 3
	}
	return result
}
