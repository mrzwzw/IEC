package IEC

import (
	"encoding/binary"
	"fmt"
)

// StartFrame 起始符
const startFrame = 0x68

// parseBigEndianUint16 转换大端Uint16
func parseBigEndianUInt16(i uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, i)
	return bytes
}

// parseLittleEndianUint16 转换小端Uint16
// func parseLittleEndianUInt16(i uint16) []byte {
// 	bytes := make([]byte, 2)
// 	binary.LittleEndian.PutUint16(bytes, i)
// 	return bytes
// }

// convertBytes 转换发送数据
func convertBytes(data []byte) []byte {
	sendData := make([]byte, 0)
	iBytes := parseBigEndianUInt16(uint16(len(data)))
	sendData = append(sendData, startFrame, iBytes[1], iBytes[1], startFrame)
	sendData = append(sendData, data...)
	return sendData
}

// 和模运算
func addMod(numbers []byte, modulus int) (result byte) {
	// numbers := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}

	intSlice := make([]int, len(numbers))
	for i, v := range numbers {
		intSlice[i] = int(v)
	}
	// 代数和
	sum := 0
	for _, num := range intSlice {
		sum += num
	}
	fmt.Println("代数和:", sum%modulus)

	// 模运算

	return byte(sum % modulus)
}
