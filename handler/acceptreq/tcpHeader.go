package accept

import "encoding/binary"

const TcpHeaderSize uint32 = 20

const FIN byte = 0x01
const SYN byte = 0x02
const RST byte = 0x04
const PSH byte = 0x08
const ACK byte = 0x10
const URG byte = 0x20

type TCPHeader struct {

	SourcePort uint16
	DestinationPort uint16
	SequenceNumber uint32
	AcknowledgementNumber uint32
	DataOffsetAndReserved byte
	HeaderLength uint32
	Flags byte
	Window uint16
	Checksum uint16
	UrgentPointer uint16
	OptionsAndPadding []byte
}

func ParseTcpHeader(buffer []byte) *TCPHeader {
	var sourcePost = binary.BigEndian.Uint16(buffer[0:2]) & 0xFFFF
	var destinationPort = binary.BigEndian.Uint16(buffer[2:4]) & 0xFFFF

	var sequenceNumber = binary.BigEndian.Uint32(buffer[4:8]) & 0xFFFFFFFF
	var acknowledgementNumber = binary.BigEndian.Uint32(buffer[8:12]) & 0xFFFFFFFF

	var dataOffsetAndReserved = buffer[12]
	var headerLength = uint32((dataOffsetAndReserved & 0xF0) >> 2)
	var flags = buffer[13]
	var window = binary.BigEndian.Uint16(buffer[14:16]) & 0xFFFF

	var checksum = binary.BigEndian.Uint16(buffer[16:18]) & 0xFFFF
	var urgentPointer = binary.BigEndian.Uint16(buffer[18:20]) & 0xFFFF

	var optionsLength = headerLength - TcpHeaderSize
	var optionsAndPadding = make([]byte,optionsLength)
	if optionsLength > 0{
		optionsAndPadding = buffer[:optionsLength]
	}

	header := &TCPHeader{
		SourcePort: sourcePost,
		DestinationPort: destinationPort,
		SequenceNumber: sequenceNumber,
		AcknowledgementNumber: acknowledgementNumber,
		DataOffsetAndReserved: dataOffsetAndReserved,
		HeaderLength: headerLength,
		Flags: flags,
		Window: window,
		Checksum: checksum,
		UrgentPointer: urgentPointer,
		OptionsAndPadding: optionsAndPadding,
	}

	return header
}

func (tcpHeader *TCPHeader) isFIN() bool{
	return (tcpHeader.Flags & FIN) == FIN
}

func (tcpHeader *TCPHeader) isSYN() bool{
	return (tcpHeader.Flags & FIN) == FIN
}

func (tcpHeader *TCPHeader) isRST() bool{
	return (tcpHeader.Flags & FIN) == FIN
}

func (tcpHeader *TCPHeader) isPSH() bool{
	return (tcpHeader.Flags & FIN) == FIN
}

func (tcpHeader *TCPHeader) isACK() bool{
	return (tcpHeader.Flags & FIN) == FIN
}
