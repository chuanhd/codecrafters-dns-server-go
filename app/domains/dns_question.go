package domains

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type DnsQuestion struct {
	Qname  string
	Qtype  uint16
	Qclass uint16
}

func (q *DnsQuestion) Encode() []byte {
	var buf bytes.Buffer

	// Split qname by '.'
	names := strings.SplitSeq(q.Qname, ".")

	for name := range names {
		buf.WriteByte(byte(len(name)))
		buf.WriteString(name)
	}
	buf.WriteByte(0)

	qTypeByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(
		qTypeByteSlice,
		q.Qtype,
	)
	buf.Write(qTypeByteSlice)

	qClassByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(qClassByteSlice, q.Qclass)
	buf.Write(qClassByteSlice)

	return buf.Bytes()
}

func DecodeQuestion(data []byte, offset int) (DnsQuestion, int, error) {
	domain, curOffset, err := parseName(data, offset)
	if err != nil {
		return DnsQuestion{}, 0, err
	}

	if curOffset+4 > len(data) {
		return DnsQuestion{}, 0, fmt.Errorf("truncated question: need 4 bytes for QTYPE+QCLASS")
	}

	qtype := binary.BigEndian.Uint16(data[curOffset : curOffset+2])
	qclass := binary.BigEndian.Uint16(data[curOffset+2 : curOffset+4])

	return DnsQuestion{
		Qname:  domain,
		Qtype:  qtype,
		Qclass: qclass,
	}, curOffset + 4, nil
}

func parseName(data []byte, offset int) (string, int, error) {
	if offset > len(data) {
		return "", 0, fmt.Errorf("Offset out of range")
	}

	domain := make([]byte, 0, 64)
	curOffset := offset

	for {
		if curOffset >= len(data) {
			return "", 0, fmt.Errorf("name parse: truncated")
		}

		curByte := data[curOffset]

		if curByte == 0x00 {
			curOffset++
			break
		}

		nameLen := int(curByte)
		curOffset++

		if curOffset+nameLen > len(data) {
			return "", 0, fmt.Errorf("name parse: truncated label")
		}

		if len(domain) > 0 {
			domain = append(domain, '.')
		}
		domain = append(domain, data[curOffset:curOffset+nameLen]...)
		curOffset += nameLen
	}

	return string(domain), curOffset, nil
}
