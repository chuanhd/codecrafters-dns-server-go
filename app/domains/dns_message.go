package domains

import "bytes"

type DnsMessage struct {
	header   DnsHeader
	question DnsQuestion
	answer   DnsAnswer
}

func CodeCraftersDnsMessage() DnsMessage {
	header := DnsHeader{
		ID:      1234,
		Flags:   0,
		QDCount: 1,
		ANCount: 1,
		NSCount: 0,
		ARCount: 0,
	}

	header.SetQR(true)

	question := DnsQuestion{
		qname:  "codecrafters.io",
		qtype:  1,
		qclass: 1,
	}

	answer := DnsAnswer{
		name:     "codecrafters.io",
		atype:    1,
		class:    1,
		ttl:      60,
		rdlength: 4,
		rdata:    "8.8.8.8",
	}

	return DnsMessage{
		header:   header,
		question: question,
		answer:   answer,
	}
}

func (m *DnsMessage) Encode() []byte {
	var outputBuff bytes.Buffer

	outputBuff.Write(m.header.Encode())
	outputBuff.Write(m.question.Encode())
	outputBuff.Write(m.answer.Encode())

	return outputBuff.Bytes()
}
