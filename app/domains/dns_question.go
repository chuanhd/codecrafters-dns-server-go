package domains

import (
	"bytes"
	"encoding/binary"
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
