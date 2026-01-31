package domains

type DnsMessage struct {
	header   DnsHeader
	question DnsQuestion
}

func EmptyDnsMessage() DnsMessage {
	header := DnsHeader{
		ID:      1234,
		Flags:   0,
		QDCount: 1,
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	header.SetQR(true)

	question := DnsQuestion{
		qname:  "codecrafters.io",
		qtype:  1,
		qclass: 1,
	}

	return DnsMessage{
		header:   header,
		question: question,
	}
}

func (m *DnsMessage) Encode() []byte {
	result := append(m.header.Encode(), m.question.Encode()...)

	return result
}
