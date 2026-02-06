package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/codecrafters-io/dns-server-starter-go/app/domains"
)

func main() {
	resolver := flag.String("resolver", "", "upstream DNS resolver address (ip:port)")
	flag.Parse()

	resolverAddr, err := net.ResolveUDPAddr("udp", *resolver)
	if err != nil {
		fmt.Printf("invalid resolver address %q: %v \n", *resolver, err)
		return
	}

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

	upstreamConn, err := net.DialUDP("udp", nil, resolverAddr)
	if err != nil {
		fmt.Println("Failed to dial upstream:", err)
		return
	}
	defer upstreamConn.Close()

	buf := make([]byte, 512)
	upstreamBuf := make([]byte, 512)

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

		questions := make([]domains.DnsQuestion, int(reqMsg.Header.QDCount))
		copy(questions, reqMsg.Question)
		var answers []domains.DnsAnswer

		for i := range questions {
			q := reqMsg.Question[i]
			fmt.Printf("Forward question %d: %s\n", i+1, q.Qname)

			upReq := buildSingleQuestionQuery(reqMsg, q)
			reqBytes := upReq.Encode()

			fmt.Printf("[your_program] >>> UPSTREAM QUERY (%d bytes)\n", len(reqBytes))
			dumpDNSBytes(reqBytes)

			// Gửi lên upstream
			if _, err := upstreamConn.Write(reqBytes); err != nil {
				fmt.Println("Failed to send upstream query:", err)
				continue
			}

			// Tránh block vô hạn
			_ = upstreamConn.SetReadDeadline(time.Now().Add(2 * time.Second))

			// Đọc response từ upstream (DialUDP => dùng Read)
			n, err := upstreamConn.Read(upstreamBuf)
			if err != nil {
				fmt.Println("Failed to read upstream response:", err)
				continue
			}

			fmt.Printf("[your_program] <<< UPSTREAM RESPONSE (%d bytes)\n", n)
			dumpDNSBytes(upstreamBuf[:n])

			upResp, err := domains.DecodeMessage(upstreamBuf[:n])
			if err != nil {
				fmt.Println("Failed to decode upstream response:", err)
				continue
			}
			fmt.Printf("upResp.Header.ANCount=%d len(upResp.Answer)=%d\n", upResp.Header.ANCount, len(upResp.Answer))

			logMessageSummary("[your_program] Decoded upstream", upResp)

			if len(upResp.Answer) > 0 {
				answers = append(answers, upResp.Answer...)
			}
		}

		respHeader := domains.DnsHeader{
			ID:      reqMsg.Header.ID,
			Flags:   0,
			QDCount: reqMsg.Header.QDCount,
			ANCount: uint16(len(answers)),
			NSCount: 0,
			ARCount: 0,
		}
		respHeader.SetResponseFlags(reqMsg.Header.Flags)

		respMsg := domains.DnsMessage{
			Header:   respHeader,
			Question: questions,
			Answer:   answers,
		}

		respBytes := respMsg.Encode()

		_, err = udpConn.WriteToUDP(respMsg.Encode(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}

		fmt.Printf("[your_program] ===== OUTGOING to %s (%d bytes) =====\n", source, len(respBytes))
		dumpDNSBytes(respBytes)

		// Decode lại chính response bytes để verify “cái mình gửi” parse ra sao
		if verify, err := domains.DecodeMessage(respBytes); err != nil {
			fmt.Println("[your_program] !!! Could not decode our own outgoing bytes:", err)

		} else {
			logMessageSummary("[your_program] Decoded outgoing(self-check)", verify)
		}

	}
}

func buildSingleQuestionQuery(orig domains.DnsMessage, q domains.DnsQuestion) domains.DnsMessage {
	flags := orig.Header.Flags &^ (1 << 15)

	h := domains.DnsHeader{
		ID:      orig.Header.ID,
		Flags:   flags,
		QDCount: 1,
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	return domains.DnsMessage{
		Header:   h,
		Question: []domains.DnsQuestion{q},
		Answer:   nil,
	}
}

// ---------- Debug helpers ----------

func dumpDNSBytes(b []byte) {
	// show first 12 bytes header quickly
	if len(b) >= 12 {
		id := uint16(b[0])<<8 | uint16(b[1])
		flags := uint16(b[2])<<8 | uint16(b[3])
		qd := uint16(b[4])<<8 | uint16(b[5])
		an := uint16(b[6])<<8 | uint16(b[7])
		ns := uint16(b[8])<<8 | uint16(b[9])
		ar := uint16(b[10])<<8 | uint16(b[11])

		fmt.Printf("[hex] ID=0x%04x FLAGS=0x%04x QD=%d AN=%d NS=%d AR=%d\n", id, flags, qd, an, ns, ar)
	}
	fmt.Println(prettyHex(b))
}

func prettyHex(b []byte) string {
	// group like hexdump but simple
	s := hex.EncodeToString(b)
	var out strings.Builder
	for i := 0; i < len(s); i += 2 {
		out.WriteString(s[i : i+2])
		if (i/2+1)%16 == 0 {
			out.WriteString("\n")
		} else {
			out.WriteString(" ")
		}
	}
	return out.String()
}

func logMessageSummary(prefix string, m domains.DnsMessage) {
	fmt.Printf("%s: id=%d flags=0x%04x qd=%d an=%d\n",
		prefix, m.Header.ID, m.Header.Flags, m.Header.QDCount, m.Header.ANCount)

	for i, q := range m.Question {
		fmt.Printf("%s: Q%d name=%q type=%d class=%d\n", prefix, i+1, q.Qname, q.Qtype, q.Qclass)
	}
	for i, a := range m.Answer {
		fmt.Printf("%s: A%d name=%q type=%d class=%d ttl=%d rdata=%v\n",
			prefix, i+1, a.Name, a.Type, a.Class, a.TTL, a.Rdata)
	}
}
