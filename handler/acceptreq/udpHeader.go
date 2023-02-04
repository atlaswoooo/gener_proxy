package accept

import "encoding/binary"

type UDPHeader struct {
	SourcePort uint16
	DestinationPort uint16

	Length uint16
	Checksum uint16
}

func ParseUDPHeader(buffer []byte) *UDPHeader {
	var sourcePort = binary.LittleEndian.Uint16(buffer[0:2]) & 0xFFFF
	var destinationPort = binary.LittleEndian.Uint16(buffer[2:4]) & 0xFFFF

	var length = binary.LittleEndian.Uint16(buffer[4:6]) & 0xFFFF
	var checksum = binary.LittleEndian.Uint16(buffer[6:8]) & 0xFFFF

	header := &UDPHeader{
		SourcePort: sourcePort,
		DestinationPort: destinationPort,
		Length: length,
		Checksum: checksum,
	}
	return header
}
