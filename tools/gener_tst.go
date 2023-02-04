package main

import (
	"encoding/binary"
	"fmt"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	UNKNOWN = Charset("UNKNOWN")
)

func main() {

	var buffer1 []byte
	var buffer2 []byte

	buffer1 = []byte{96, 1, 3, 0, 0, 36, 0, 1, 254, 128, 0}
	buffer2 = buffer1[2:6]

	totalLength := binary.LittleEndian.Uint32(buffer2) & 0xFFFF
	fmt.Printf("totalLength:%+v\n", totalLength)

}
