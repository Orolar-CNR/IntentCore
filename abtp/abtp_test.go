package abtp

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// AbtpHeader represents the ABTP frame header
type AbtpHeader struct {
	Magic   uint32
	Version uint16
	Length  uint16
}

// Simulated simple frame parsing to test Go overhead and logic similar to XDP.
// In a real environment, XDP runs in kernel space before reaching this Go code.
// This benchmark serves as a Go-level validation/parser latency test.
func ParseABTPFrame(data []byte) (AbtpHeader, error) {
	if len(data) < 8 {
		return AbtpHeader{}, bytes.ErrTooLarge // Using as a generic error
	}
	var hdr AbtpHeader
	hdr.Magic = binary.BigEndian.Uint32(data[0:4])
	hdr.Version = binary.BigEndian.Uint16(data[4:6])
	hdr.Length = binary.BigEndian.Uint16(data[6:8])

	if hdr.Magic != 0x41425450 {
		return hdr, bytes.ErrTooLarge
	}
	if hdr.Version != 1 {
		return hdr, bytes.ErrTooLarge
	}
	if hdr.Length < 8 {
		return hdr, bytes.ErrTooLarge
	}

	return hdr, nil
}

func BenchmarkParseABTPFrame_Valid(b *testing.B) {
	// Construct a valid ABTP frame header
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint32(0x41425450)) // Magic "ABTP"
	binary.Write(&buf, binary.BigEndian, uint16(1))          // Version 1
	binary.Write(&buf, binary.BigEndian, uint16(8))          // Length 8

	data := buf.Bytes()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ParseABTPFrame(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseABTPFrame_InvalidMagic(b *testing.B) {
	// Construct an invalid ABTP frame header
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint32(0x00000000)) // Invalid Magic
	binary.Write(&buf, binary.BigEndian, uint16(1))
	binary.Write(&buf, binary.BigEndian, uint16(8))

	data := buf.Bytes()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseABTPFrame(data) // Expecting an error
	}
}
