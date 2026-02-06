package domains

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/dns-server-starter-go/app/utils"
)

type DnsAnswer struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	Rdata []byte
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

	rData := make([]byte, 4)
	copy(rData[:], q.Rdata)

	rlengthByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(rlengthByteSlice, uint16(len(rData)))
	buf.Write(rlengthByteSlice)

	buf.Write(rData[:])

	return buf.Bytes()
}

func DecodeAnswer(data []byte, offset int) (DnsAnswer, int, error) {
	name, curOffset, err := utils.DecodeName(data, offset)
	if err != nil {
		return DnsAnswer{}, offset, err
	}

	type_ := binary.BigEndian.Uint16(data[curOffset : curOffset+2])
	class := binary.BigEndian.Uint16(data[curOffset+2 : curOffset+4])
	ttl := binary.BigEndian.Uint32(data[curOffset+4 : curOffset+8])
	rlength := binary.BigEndian.Uint16(data[curOffset+8 : curOffset+10])

	curOffset = curOffset + 10

	rdata := make([]byte, int(rlength))
	copy(rdata, data[curOffset:curOffset+int(rlength)])

	curOffset += int(rlength)

	return DnsAnswer{
		Name:  name,
		Type:  type_,
		Class: class,
		TTL:   ttl,
		Rdata: rdata,
	}, curOffset, nil
}
