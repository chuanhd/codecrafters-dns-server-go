package domains

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type DnsQuestion struct {
	qname  string
	qtype  uint16
	qclass uint16
}

func (q *DnsQuestion) Encode() []byte {
	var buf bytes.Buffer

	// Split qname by '.'
	names := strings.SplitSeq(q.qname, ".")

	for name := range names {
		buf.WriteByte(byte(len(name)))
		buf.WriteString(name)
	}
	buf.WriteByte(0)

	qTypeByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(
		qTypeByteSlice,
		q.qtype,
	)
	buf.Write(qTypeByteSlice)

	qClassByteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(qClassByteSlice, q.qclass)
	buf.Write(qClassByteSlice)

	return buf.Bytes()
}
