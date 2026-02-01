package domains

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type DnsAnswer struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDlength uint16
	Rdata    string
}

func encodeRData(ipStr string) ([4]byte, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return [4]byte{}, fmt.Errorf("invalid ip: %s", ipStr)
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return [4]byte{}, fmt.Errorf("not ipv4: %s", ipStr)
	}
	var r [4]byte
	copy(r[:], ip4) // r = [a b c d]
	return r, nil
}

func (q *DnsAnswer) Encode() []byte {
	var buf bytes.Buffer

	// Split qname by '.'
	names := strings.SplitSeq(q.Name, ".")

	for name := range names {
		buf.WriteByte(byte(len(name)))
		buf.WriteString(name)
	}
	buf.WriteByte(0)

	typeByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(
		typeByteSlice,
		q.Type,
	)
	buf.Write(typeByteSlice)

	qClassByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(qClassByteSlice, q.Class)
	buf.Write(qClassByteSlice)

	ttlByteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(ttlByteSlice, q.TTL)
	buf.Write(ttlByteSlice)

	rlengthByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(rlengthByteSlice, q.RDlength)
	buf.Write(rlengthByteSlice)

	rDataInByte, _ := encodeRData(q.Rdata)
	buf.Write(rDataInByte[:])

	return buf.Bytes()
}
