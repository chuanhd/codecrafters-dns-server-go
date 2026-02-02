package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/domains"
)

// Ensures gofmt doesn't remove the "net" import in stage 1 (feel free to remove this!)
var _ = net.ListenUDP

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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
			QDCount: 1,
			ANCount: 1,
			NSCount: 0,
			ARCount: 0,
		}
		respHeader.SetResponseFlags(reqMsg.Header.Flags)

		question := domains.DnsQuestion{
			Qname:  reqMsg.Question.Qname,
			Qtype:  1,
			Qclass: 1,
		}

		answer := domains.DnsAnswer{
			Name:     reqMsg.Question.Qname,
			Type:     1,
			Class:    1,
			TTL:      60,
			RDlength: 4,
			Rdata:    "8.8.8.8",
		}

		respMsg := domains.DnsMessage{
			Header:   respHeader,
			Question: question,
			Answer:   answer,
		}

		_, err = udpConn.WriteToUDP(respMsg.Encode(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
