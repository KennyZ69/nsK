package main

type NetPacket interface {
	// Marshall(payload []byte) ([]byte, error)
	// UnMarshall(p []byte) (*NetPacket, error)
	Size() int
}

type SimPacket struct {
	Source  Node
	Dest    Node
	Ack     bool // whether the packet was acknowledged
	Payload []byte
}

type IPPacket struct {
	Version            uint8  // 4bit
	HeaderSize         uint8  // 4 bit
	ServiceType        uint8  // 8 bit
	TotalSize          uint16 // 16 bit
	Identifier         uint16 // 16 bit
	Flags              uint8  // 3 bit
	FragmentOffset     uint16 // 13 bit
	TTL                uint8  // 8 bit
	Protocol           uint8  // 8 bit
	HeaderChecksum     uint16 // 16 bit
	SourceAddress      uint32 // 32 bit
	DestinationAddress uint32 // 32 bit
	// Options            []byte // vary
	Data NetPacket
}

func (p *IPPacket) Size() int {
	return int(p.TotalSize)
}
