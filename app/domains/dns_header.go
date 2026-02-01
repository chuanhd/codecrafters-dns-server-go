package domains

import (
	"encoding/binary"
	"fmt"
)

type DnsHeader struct {
	ID      uint16 // 16 bits
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

const (
	flagQR     uint16 = 1 << 15      // 1 bit
	maskOpCode uint16 = 0b1111 << 11 // 4 bits (14..11)
	flagAA     uint16 = 1 << 10
	flagTC     uint16 = 1 << 9
	flagRD     uint16 = 1 << 8
	flagRA     uint16 = 1 << 7
	maskZ      uint16 = 0b111 << 4 // 3 bits (6..4)
	maskRCODE  uint16 = 0b1111
)

func (h *DnsHeader) SetQR(v bool) {
	if v {
		h.Flags |= flagQR
	} else {
		h.Flags &^= flagQR
	}
}

func (h *DnsHeader) Opcode() uint16 {
	return (h.Flags & maskOpCode) >> 11
}
func (h *DnsHeader) SetOpcode(op uint16) {
	h.Flags &^= maskOpCode
	h.Flags |= (op & 0b1111) << 11
}

func (h *DnsHeader) RCode() uint16 {
	return h.Flags & maskRCODE
}
func (h *DnsHeader) SetRCode(rc uint16) {
	h.Flags &^= maskRCODE
	h.Flags |= (rc & 0b1111)
}

func (h *DnsHeader) SetResponseFlags(requestFlags uint16) {
	opCode := requestFlags & maskOpCode
	rd := requestFlags & flagRD

	var rcode uint16 = 0
	if opCode != 0 {
		rcode = 4
	}

	h.Flags = flagQR | opCode | rd | rcode
}

func (h *DnsHeader) Encode() []byte {
	b := make([]byte, 12)

	binary.BigEndian.PutUint16(b[0:2], h.ID)
	binary.BigEndian.PutUint16(b[2:4], h.Flags)
	binary.BigEndian.PutUint16(b[4:6], h.QDCount)
	binary.BigEndian.PutUint16(b[6:8], h.ANCount)
	binary.BigEndian.PutUint16(b[8:10], h.NSCount)
	binary.BigEndian.PutUint16(b[10:12], h.ARCount)

	return b
}

func DecodeHeader(b []byte) (DnsHeader, error) {
	if len(b) < 12 {
		return DnsHeader{}, fmt.Errorf("Insufficient data of header")
	}

	h := DnsHeader{
		ID:      binary.BigEndian.Uint16(b[0:2]),
		Flags:   binary.BigEndian.Uint16(b[2:4]),
		QDCount: binary.BigEndian.Uint16(b[4:6]),
		ANCount: binary.BigEndian.Uint16(b[6:8]),
		NSCount: binary.BigEndian.Uint16(b[8:10]),
		ARCount: binary.BigEndian.Uint16(b[10:12]),
	}
	return h, nil
}
