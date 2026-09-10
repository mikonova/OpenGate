package payloads

import (
	"encoding/binary"
	"time"
)

/*
This is a simple implementation of a byte protocol.
1st byte is a packet type from types subpackage,
8 bytes after are a timestamp in form of Time.Time, converted to uint64.
IsOutgoing is a boolean, which takes 1 bit, but aligns to 1 byte.
MessageLength is the length of a string, int32 in form of 4 bytes.
Message is aa variable string length. Pay attention, that the byte sequence aligns to 4 bytes by padding with zeros ao the end,
which is calculated in Decode() function. Please use the constructor function for the payloads.Packet type instead of manual creation
*/
type Packet struct {
	PacketType    byte      // 1 byte
	Timestamp     time.Time // 8 bytes uint64
	IsOutgoing    bool      // 1 byte
	MessageLength int32     // 4 bytes
	Message       string    // variable
}

func (p Packet) Encode() (bytePacket []byte) {
	bytePacket = make([]byte, 0)
	bytePacket = append(bytePacket, p.PacketType)
	binary.BigEndian.AppendUint64(bytePacket, uint64(p.Timestamp.Unix()))
	if p.IsOutgoing == true {
		bytePacket = append(bytePacket, 0b00000001)
	} else {
		bytePacket = append(bytePacket, 0b00000000)
	}
	binary.BigEndian.AppendUint32(bytePacket, uint32(p.MessageLength))
	bytePacket = append(bytePacket, []byte(p.Message)...)
	if p.MessageLength%4 != 0 {
		bytePacket = append(bytePacket, make([]byte, int(4-p.MessageLength%4))...)
	}
	return
}

func (p Packet) Decode(bytePacket []byte) {
	p.PacketType = bytePacket[1]
	unixSeconds := binary.BigEndian.Uint64(bytePacket[2:10])
	p.Timestamp = time.Unix(int64(unixSeconds), 0)
	if bytePacket[10] == 0x00 {
		p.IsOutgoing = true
	} else {
		p.IsOutgoing = false
	}
	p.MessageLength = int32(binary.BigEndian.Uint32(bytePacket[11:15]))
	p.Message = string(bytePacket[15 : p.MessageLength-1])
}

func NewPacket(pType byte, timestamp time.Time, isOutgoing bool, message string) (p Packet) {
	p = Packet{
		PacketType:    pType,
		Timestamp:     timestamp,
		IsOutgoing:    isOutgoing,
		MessageLength: int32(len(message)),
		Message:       message,
	}
	return
}
