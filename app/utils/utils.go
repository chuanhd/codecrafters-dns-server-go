package utils

import (
	"encoding/binary"
	"fmt"
)

func DecodeName(data []byte, offset int) (string, int, error) {
	if offset > len(data) {
		return "", 0, fmt.Errorf("Offset out of range")
	}

	domain := make([]byte, 0, 64)
	curOffset := offset
	jumped := false
	afterJumpOffset := 0 // Store the offset after jumping

	for {
		if curOffset >= len(data) {
			return "", 0, fmt.Errorf("name parse: truncated")
		}

		curByte := data[curOffset]

		// Compression pointer
		if curByte&0xC0 == 0xC0 {
			if curOffset+2 > len(data) {
				return "", 0, fmt.Errorf("name parse: truncated compression pointer")
			}
			compressedOffset := int(binary.BigEndian.Uint16(data[curOffset:curOffset+2]) & 0x3FFF)
			if compressedOffset > len(data) {
				return "", 0, fmt.Errorf("name parse: invalid compression pointer")
			}
			if !jumped {
				jumped = true
				afterJumpOffset = curOffset + 2
			}
			curOffset = compressedOffset
			continue
		}

		if curByte == 0x00 {
			curOffset++
			break
		}

		nameLen := int(curByte)
		curOffset++

		if curOffset+nameLen > len(data) {
			return "", 0, fmt.Errorf("name parse: truncated label")
		}

		if len(domain) > 0 {
			domain = append(domain, '.')
		}
		domain = append(domain, data[curOffset:curOffset+nameLen]...)
		curOffset += nameLen
	}

	name := string(domain)
	if jumped {
		return name, afterJumpOffset, nil
	}

	return name, curOffset, nil
}
