package main

import (
	"fmt"
)

const (
	BATCH_SIZE       = 3
	MANAGE_BYTE byte = 255
	DICT_SIZE        = 255*256 + 255
)

func incValue(value []byte) []byte {
	if value[1] < MANAGE_BYTE {
		value[1]++
	} else {
		value[0]++
		value[1] = 0
	}
	return value
}

func getBufferForArchive(buffer [][]byte) ([]byte, [][]byte, []byte) {
	var result []byte
	var dictSlice [][]byte
	dict := make(map[string][]byte)
	value := []byte{0, 0}
	for index, batch := range buffer {
		if len(batch)-BATCH_SIZE < 0 {
			result = appendManageBytesWithBatch(result, batch)
			break
		}
		key := string(batch)
		if _, ok := dict[key]; !ok {
			if len(dict) < DICT_SIZE {
				dict[key] = []byte{value[0], value[1]}
				dictSlice = append(dictSlice, batch)
				if index != 0 {
					result = append(result, value[0], value[1])
				}
			} else {
				result = appendManageBytesWithBatch(result, batch)
			}
			value = incValue(value)
		} else {
			result = append(result, dict[key][0], dict[key][1])
		}
	}
	dictSize := []byte{byte(len(dict) / 256), byte(len(dict) % 256)}
	return result, dictSlice, dictSize
}

func appendManageBytesWithBatch(result []byte, batch []byte) []byte {
	result = append(result, MANAGE_BYTE, MANAGE_BYTE)
	for _, val := range batch {
		result = append(result, val)
	}
	return result
}

func decodeBuffer(buffer []byte) [][]byte {
	if !validateExtension(buffer) {
		return nil
	}
	var decodedBuffer [][]byte
	dictLen := getDictLen(buffer)
	dict, dictSlice := getDict(buffer, dictLen)
	decodedBuffer = append(decodedBuffer, []byte{dictSlice[0]})
	bufferIndex := len(dictSlice) + 5
	bufferLen := len(buffer)
	for i := bufferIndex; i < bufferLen; i += 2 {
		if buffer[i] == MANAGE_BYTE && buffer[i+1] == MANAGE_BYTE {
			if bufferLen-(i+2) < BATCH_SIZE {
				decodedBuffer = append(decodedBuffer, []byte{buffer[i+2], buffer[i+3]})
				continue
			}
			decodedBuffer = append(decodedBuffer, []byte{buffer[i+3], buffer[i+4]})
			i += 3
		} else {
			x := string([]byte{buffer[i], buffer[i+1]})
			decodedBuffer = append(decodedBuffer, dict[x])
		}
	}
	return decodedBuffer
}

func getDict(buffer []byte, dictLen int) (map[string][]byte, []byte) {
	dict := make(map[string][]byte)
	var dictSlice []byte
	value := []byte{0, 0}
	for i := 5; i < (dictLen*3)+5; i += 3 {
		batch := make([]byte, 0, BATCH_SIZE)
		batch = append(batch, buffer[i], buffer[i+1], buffer[i+2])
		str := string([]byte{value[0], value[1]})
		dict[str] = batch
		value = incValue(value)
		dictSlice = append(dictSlice, buffer[i], buffer[i+1], buffer[i+2])
	}
	return dict, dictSlice
}

func validateExtension(buffer []byte) bool {
	var extension []byte
	for i := 0; i < BATCH_SIZE; i++ {
		extension = append(extension, buffer[i])
	}
	if string(extension) != ARCHIVE_EXTENSION {
		fmt.Println("INVALID FORMAT")
		return false
	}

	return true
}

func getDictLen(buffer []byte) int {
	div := int(buffer[3])
	mod := int(buffer[4])
	return div*256 + mod
}
