package domains

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/codecrafters-io/dns-server-starter-go/app/utils"
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
	domain, curOffset, err := utils.DecodeName(data, offset)
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
