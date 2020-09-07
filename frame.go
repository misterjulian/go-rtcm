package rtcm3

import (
	"bufio"
	"encoding/binary"
	"errors"

	"github.com/bamiaux/iobit"

	"github.com/misterjulian/go-rtcm/messages"
)

const (
	FramePreamble byte = 0xD3

	ErrInvalidPreamble = "invalid preamble"
	ErrInvalidCRC      = "invalid crc"
)

type Frame struct {
	Preamble uint8
	Reserved uint8
	Length   uint16
	Payload  []byte
	Crc      uint32
}

// EncapsulateByteArray lazily wraps any byte array in an RTCM3 Frame
func EncapsulateByteArray(data []byte) (frame *Frame) {
	frame = &Frame{
		Preamble: FramePreamble,
		Reserved: 0,
		Length:   uint16(len(data)),
		Payload:  data,
		Crc:      uint32(0),
	}
	frame.Crc = Crc24q(frame.Serialize()[:len(data)+3])
	return frame
}

func EncapsulateMessage(msg messages.Message) (frame *Frame) {
	return EncapsulateByteArray(msg.Serialize())
}

func (f *Frame) MessageNumber() uint16 {
	return binary.BigEndian.Uint16(f.Payload[0:2]) >> 4
}

func (f *Frame) Serialize() []byte {
	data := make([]byte, f.Length+6)
	w := iobit.NewWriter(data)
	w.PutUint8(8, f.Preamble)
	w.PutUint8(6, f.Reserved)
	w.PutUint16(10, f.Length)
	w.Write(f.Payload)
	w.PutUint32(24, f.Crc)
	w.Flush()
	return data
}

func DeserializeFrame(reader *bufio.Reader) (frame *Frame, err error) {
	// Only reads first byte from reader if Preamble or CRC are incorrect
	// Unfortunatly can't construct anything that will read bits (like iobit) until we have a byte array
	preamble, err := reader.ReadByte()
	if err != nil {
		return frame, err
	}
	if preamble != FramePreamble {
		return frame, errors.New(ErrInvalidPreamble)
	}

	header, err := reader.Peek(2)
	if err != nil {
		return frame, err
	}

	reserved := uint8(header[0]) & 0xFC
	length := binary.BigEndian.Uint16(header) & 0x3FF
	data, err := reader.Peek(int(length + 5))
	if err != nil {
		return frame, err
	}

	data = append([]byte{preamble}, data...)
	crc := binary.BigEndian.Uint32(data[len(data)-4:]) & 0xFFFFFF

	frame = &Frame{
		Preamble: preamble,
		Reserved: reserved,
		Length:   length,
		Payload:  data[3 : len(data)-3],
		Crc:      crc,
	}

	if Crc24q(data[:len(data)-3]) != frame.Crc {
		return frame, errors.New(ErrInvalidCRC)
	}

	_, _ = reader.Discard(len(data) - 1)
	return frame, nil
}
