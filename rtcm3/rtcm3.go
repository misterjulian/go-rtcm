package rtcm3

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/bamiaux/iobit"
	"io"
	"time"
)

type Message interface {
	Serialize() []byte
	Number() uint16
}

func DeserializeMessage(payload []byte) (msg Message) {
	messageNumber := binary.BigEndian.Uint16(payload[0:2]) >> 4

	switch int(messageNumber) {
	case 1001:
		return DeserializeMessage1001(payload)
	case 1002:
		return DeserializeMessage1002(payload)
	case 1003:
		return DeserializeMessage1003(payload)
	case 1004:
		return DeserializeMessage1004(payload)
	case 1005:
		return DeserializeMessage1005(payload)
	case 1006:
		return DeserializeMessage1006(payload)
	case 1007:
		return DeserializeMessage1007(payload)
	case 1008:
		return DeserializeMessage1008(payload)
	case 1009:
		return DeserializeMessage1009(payload)
	case 1010:
		return DeserializeMessage1010(payload)
	case 1011:
		return DeserializeMessage1011(payload)
	case 1012:
		return DeserializeMessage1012(payload)
	case 1013:
		return DeserializeMessage1013(payload)
	case 1014:
		return DeserializeMessage1014(payload)
	case 1015:
		return DeserializeMessage1015(payload)
	case 1016:
		return DeserializeMessage1016(payload)
	case 1019:
		return DeserializeMessage1019(payload)
	case 1020:
		return DeserializeMessage1020(payload)
	case 1029:
		return DeserializeMessage1029(payload)
	case 1030:
		return DeserializeMessage1030(payload)
	case 1031:
		return DeserializeMessage1031(payload)
	case 1032:
		return DeserializeMessage1032(payload)
	case 1033:
		return DeserializeMessage1033(payload)
	case 1071:
		return DeserializeMessage1071(payload)
	case 1072:
		return DeserializeMessage1072(payload)
	case 1073:
		return DeserializeMessage1073(payload)
	case 1074:
		return DeserializeMessage1074(payload)
	case 1075:
		return DeserializeMessage1075(payload)
	case 1076:
		return DeserializeMessage1076(payload)
	case 1077:
		return DeserializeMessage1077(payload)
	case 1081:
		return DeserializeMessage1081(payload)
	case 1082:
		return DeserializeMessage1082(payload)
	case 1083:
		return DeserializeMessage1083(payload)
	case 1084:
		return DeserializeMessage1084(payload)
	case 1085:
		return DeserializeMessage1085(payload)
	case 1086:
		return DeserializeMessage1086(payload)
	case 1087:
		return DeserializeMessage1087(payload)
	case 1091:
		return DeserializeMessage1091(payload)
	case 1092:
		return DeserializeMessage1092(payload)
	case 1093:
		return DeserializeMessage1093(payload)
	case 1094:
		return DeserializeMessage1094(payload)
	case 1095:
		return DeserializeMessage1095(payload)
	case 1096:
		return DeserializeMessage1096(payload)
	case 1097:
		return DeserializeMessage1097(payload)
	case 1101:
		return DeserializeMessage1101(payload)
	case 1102:
		return DeserializeMessage1102(payload)
	case 1103:
		return DeserializeMessage1103(payload)
	case 1104:
		return DeserializeMessage1104(payload)
	case 1105:
		return DeserializeMessage1105(payload)
	case 1106:
		return DeserializeMessage1106(payload)
	case 1107:
		return DeserializeMessage1107(payload)
	case 1111:
		return DeserializeMessage1111(payload)
	case 1112:
		return DeserializeMessage1112(payload)
	case 1113:
		return DeserializeMessage1113(payload)
	case 1114:
		return DeserializeMessage1114(payload)
	case 1115:
		return DeserializeMessage1115(payload)
	case 1116:
		return DeserializeMessage1116(payload)
	case 1117:
		return DeserializeMessage1117(payload)
	case 1121:
		return DeserializeMessage1121(payload)
	case 1122:
		return DeserializeMessage1122(payload)
	case 1123:
		return DeserializeMessage1123(payload)
	case 1124:
		return DeserializeMessage1124(payload)
	case 1125:
		return DeserializeMessage1125(payload)
	case 1126:
		return DeserializeMessage1126(payload)
	case 1127:
		return DeserializeMessage1127(payload)
	case 1230:
		return DeserializeMessage1230(payload)
	default:
		return MessageUnknown{payload}
	}
}

type MessageUnknown struct {
	Payload []byte
}

func (msg MessageUnknown) Serialize() []byte {
	return msg.Payload
}

func (msg MessageUnknown) Number() (msgNumber uint16) {
	return binary.BigEndian.Uint16(msg.Payload[0:4]) >> 4
}

type Observable interface {
	Time() time.Time
}

var FramePreamble byte = 0xD3

type Frame struct {
	Preamble uint8
	Reserved uint8
	Length   uint16
	Payload  []byte
	Crc      uint32
}

func EncapsulateMessage(msg Message) (frame Frame) {
	data := msg.Serialize()
	frame = Frame{
		Preamble: FramePreamble,
		Reserved: 0,
		Length:   uint16(len(data)),
		Payload:  data,
		Crc:      uint32(0),
	}
	frame.Crc = Crc24q(frame.Serialize()[:len(data)+3])
	return frame
}

func (frame Frame) Serialize() []byte {
	data := make([]byte, frame.Length+6)
	w := iobit.NewWriter(data)
	w.PutUint8(8, frame.Preamble)
	w.PutUint8(6, frame.Reserved)
	w.PutUint16(10, frame.Length)
	w.Write(frame.Payload)
	w.PutUint32(24, frame.Crc)
	w.Flush()
	return data
}

func DeserializeFrame(reader *bufio.Reader) (frame Frame, err error) {
	// Only reads first byte from reader if Preamble or CRC are incorrect
	// Unfortunatly can't construct anything that will read bits (like iobit) until we have a byte array
	preamble, err := reader.ReadByte()
	if err != nil {
		return frame, err
	}
	if preamble != FramePreamble {
		return frame, errors.New("Invalid Preamble")
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

	frame = Frame{
		Preamble: preamble,
		Reserved: reserved,
		Length:   length,
		Payload:  data[3 : len(data)-3],
		Crc:      crc,
	}

	if Crc24q(data[:len(data)-3]) != frame.Crc {
		return frame, errors.New("CRC Failed")
	}

	reader.Discard(len(data) - 1)
	return frame, nil
}

type Scanner struct {
	Reader *bufio.Reader
}

func NewScanner(r io.Reader) Scanner {
	return Scanner{bufio.NewReader(r)}
}

func (scanner Scanner) Next() (message Message, err error) {
	for {
		frame, err := DeserializeFrame(scanner.Reader)
		if err != nil {
			if err.Error() == "Invalid Preamble" || err.Error() == "CRC Failed" {
				continue
			}
			return nil, err
		}
		return DeserializeMessage(frame.Payload), err // probably have frame.Message() return err
	}
}
