package domains

import "bytes"

type DnsMessage struct {
	Header   DnsHeader
	Question DnsQuestion
	Answer   DnsAnswer
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
		Qname:  "codecrafters.io",
		Qtype:  1,
		Qclass: 1,
	}

	answer := DnsAnswer{
		Name:     "codecrafters.io",
		Type:     1,
		Class:    1,
		TTL:      60,
		RDlength: 4,
		Rdata:    "8.8.8.8",
	}

	return DnsMessage{
		Header:   header,
		Question: question,
		Answer:   answer,
	}
}

func (m *DnsMessage) Encode() []byte {
	var outputBuff bytes.Buffer

	outputBuff.Write(m.Header.Encode())
	outputBuff.Write(m.Question.Encode())
	outputBuff.Write(m.Answer.Encode())

	return outputBuff.Bytes()
}

func DecodeMessage(data []byte) (DnsMessage, error) {
	headerInBytes := data[:12]

	header, err := DecodeHeader(headerInBytes)

	if err != nil {
		return DnsMessage{}, err
	}

	header.SetQR(true)
	header.ANCount = 1

	question := DnsQuestion{
		Qname:  "codecrafters.io",
		Qtype:  1,
		Qclass: 1,
	}

	answer := DnsAnswer{
		Name:     "codecrafters.io",
		Type:     1,
		Class:    1,
		TTL:      60,
		RDlength: 4,
		Rdata:    "8.8.8.8",
	}

	return DnsMessage{
		Header:   header,
		Question: question,
		Answer:   answer,
	}, nil
}
