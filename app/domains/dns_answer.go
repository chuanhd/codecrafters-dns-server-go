package domains

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type DnsAnswer struct {
	name     string
	atype    uint16
	class    uint16
	ttl      uint32
	rdlength uint16
	rdata    string
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
	names := strings.SplitSeq(q.name, ".")

	for name := range names {
		buf.WriteByte(byte(len(name)))
		buf.WriteString(name)
	}
	buf.WriteByte(0)

	typeByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(
		typeByteSlice,
		q.atype,
	)
	buf.Write(typeByteSlice)

	qClassByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(qClassByteSlice, q.class)
	buf.Write(qClassByteSlice)

	ttlByteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(ttlByteSlice, q.ttl)
	buf.Write(ttlByteSlice)

	rlengthByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(rlengthByteSlice, q.rdlength)
	buf.Write(rlengthByteSlice)

	rDataInByte, _ := encodeRData(q.rdata)
	buf.Write(rDataInByte[:])

	return buf.Bytes()
}
