package rtcm3

import (
	"bufio"
	"io"

	"github.com/misterjulian/go-rtcm/messages"
)

type (
	Scanner interface {
		Scan() (message messages.Message, err error)
		ScanFrame() (frame *Frame, err error)
	}

	scanner struct {
		Reader *bufio.Reader
	}
)

// New returns a Rtcm interface for getting next rtcm v3 message
func NewScanner(r io.Reader) Scanner {
	return &scanner{Reader: bufio.NewReader(r)}
}

func (s *scanner) Scan() (message messages.Message, err error) {
	frame, err := s.ScanFrame()
	if err != nil {
		return nil, err
	}
	return messages.DeserializeMessage(frame.Payload), err // probably have DeserializeMessage return err
}

func (s *scanner) ScanFrame() (frame *Frame, err error) {
	for {
		frame, err := DeserializeFrame(s.Reader)
		if err != nil {
			if err.Error() == ErrInvalidPreamble || err.Error() == ErrInvalidCRC {
				// Continue reading from next byte if a valid Frame was not found
				//TODO: return byte array of skipped bytes
				continue
			}
		}
		return frame, err
	}
}
