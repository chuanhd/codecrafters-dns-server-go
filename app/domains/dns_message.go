package domains

import (
	"bytes"
)

type DnsMessage struct {
	Header   DnsHeader
	Question []DnsQuestion
	Answer   []DnsAnswer
}

func (m *DnsMessage) Encode() []byte {
	var outputBuff bytes.Buffer

	outputBuff.Write(m.Header.Encode())
	for _, question := range m.Question {
		outputBuff.Write(question.Encode())
	}
	for _, answer := range m.Answer {
		outputBuff.Write(answer.Encode())
	}

	return outputBuff.Bytes()
}

func DecodeMessage(data []byte) (DnsMessage, error) {
	const headerLength = 12

	headerInBytes := data[:headerLength]

	header, err := DecodeHeader(headerInBytes)

	if err != nil {
		return DnsMessage{}, err
	}

	offset := headerLength
	questions := make([]DnsQuestion, 0)
	for i := 0; i < int(header.QDCount); i++ {
		question, next, err := DecodeQuestion(data, offset)
		offset = next

		if err != nil {
			continue
		}

		questions = append(questions, question)
	}

	answers := make([]DnsAnswer, 0, int(header.ANCount))
	for i := 0; i < int(header.ANCount); i++ {
		a, next, err := DecodeAnswer(data, offset)
		if err != nil {
			return DnsMessage{}, err
		}
		answers = append(answers, a)
		offset = next
	}

	return DnsMessage{
		Header:   header,
		Question: questions,
		Answer:   answers,
	}, nil
}
