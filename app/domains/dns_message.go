package domains

type DnsMessage struct {
	header DnsHeader
}

func EmptyDnsMessage() DnsMessage {
	header := DnsHeader{
		ID:      1234,
		Flags:   0,
		QDCount: 0,
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	header.SetQR(true)

	return DnsMessage{
		header: header,
	}
}

func (m *DnsMessage) Encode() []byte {
	return m.header.Encode()
}
