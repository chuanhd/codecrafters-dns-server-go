package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/domains"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		reqMsg, err := domains.DecodeMessage(buf[:size])
		if err != nil {
			fmt.Println("Error decoding message: ", err)
			continue
		}

		respHeader := domains.DnsHeader{
			ID:      reqMsg.Header.ID,
			Flags:   0,
			QDCount: reqMsg.Header.QDCount,
			ANCount: reqMsg.Header.QDCount,
			NSCount: 0,
			ARCount: 0,
		}
		respHeader.SetResponseFlags(reqMsg.Header.Flags)

		questions := make([]domains.DnsQuestion, int(reqMsg.Header.QDCount))
		answers := make([]domains.DnsAnswer, int(reqMsg.Header.QDCount))

		for i := range questions {
			fmt.Printf("Question %d: %s\n", i+1, reqMsg.Question[i].Qname)
			questions[i] = domains.DnsQuestion{
				Qname:  reqMsg.Question[i].Qname,
				Qtype:  1,
				Qclass: 1,
			}

			answers[i] = domains.DnsAnswer{
				Name:  reqMsg.Question[i].Qname,
				Type:  1,
				Class: 1,
				TTL:   60,
				Rdata: "8.8.8.8",
			}
		}

		respMsg := domains.DnsMessage{
			Header:   respHeader,
			Question: questions,
			Answer:   answers,
		}

		_, err = udpConn.WriteToUDP(respMsg.Encode(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
