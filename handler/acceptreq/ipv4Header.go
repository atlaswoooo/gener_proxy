package accept

import (
	"encoding/binary"
	"fmt"
)

type IPV4Header struct {
	 VersionAndIHL byte
	 Version byte
	 IHL byte
	 HeaderLength byte
	 TypeOfService uint16
	 TotalLength uint16
	 IdentificationAndFlagsAndFragmentOffset uint32
	 TTL byte
	 ProtocolNum byte
	 Protocol string
	 HeaderChecksum uint16
	 SourceAddress string
	 DestinationAddress string
}

func ParseIPV4Header(buffer []byte) *IPV4Header {

	var versionAndIHL = buffer[0]
	var version = versionAndIHL >> 4
	var IHL = versionAndIHL & 0x0F
	var headerLength = IHL << 2
	var typeOfService = uint16(buffer[1] & 0xFF)
	var totalLength = binary.LittleEndian.Uint16(buffer[2:4]) & 0xFFFF
	var identificationAndFlagsAndFragmentOffset = binary.LittleEndian.Uint32(buffer[4:8]) & 0xFFFFFFFF
	var TTL = buffer[8] & 0x0F
	var protocolNum = buffer[9] & 0x0F
	var protocol string
	if protocolNum == 6 {
		protocol  = "TCP"
	} else if protocolNum == 17 {
		protocol = "UDP"
	}
	var headerChecksum = binary.LittleEndian.Uint16(buffer[10:12]) & 0xFFFF
	var sourceAddress = fmt.Sprintf("%d.%d.%d.%d", buffer[12], buffer[13], buffer[14], buffer[15])
	var destinationAddress = fmt.Sprintf("%d.%d.%d.%d", buffer[16], buffer[17], buffer[18], buffer[19])

	header := &IPV4Header{
		VersionAndIHL:versionAndIHL,
		Version: version,
		IHL: IHL,
		HeaderLength: headerLength,
		TypeOfService: typeOfService,
		TotalLength: totalLength,
		IdentificationAndFlagsAndFragmentOffset: identificationAndFlagsAndFragmentOffset,
		TTL: TTL,
		ProtocolNum: protocolNum,
		Protocol: protocol,
		HeaderChecksum: headerChecksum,
		SourceAddress: sourceAddress,
		DestinationAddress: destinationAddress,
	}

	return header
}

func (ipHeader *IPV4Header)Print()  {
	fmt.Printf("destinationIp:%v\n",ipHeader.DestinationAddress)
	fmt.Printf("sourceIp:%v\n",ipHeader.SourceAddress)
	fmt.Printf("protocol:%v\n",ipHeader.Protocol)
	fmt.Println()
}

