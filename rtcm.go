package rtcm3

import (
	"bufio"
	"context"
	"io"

	"github.com/misterjulian/go-rtcm/messages"
)

type (
	Rtcm interface {
		Next() (message messages.Message, err error)
		NextFrame() (frame *Frame, err error)
	}

	rtcm struct {
		ctx    context.Context
		Reader *bufio.Reader
	}
)

// New returns a Rtcm interface for getting next rtcm v3 message
func New(ctx context.Context, r io.Reader) Rtcm {
	return &rtcm{ctx: ctx, Reader: bufio.NewReader(r)}
}

func (s *rtcm) Next() (message messages.Message, err error) {
	frame, err := s.NextFrame()
	if err != nil {
		return nil, err
	}
	return messages.DeserializeMessage(frame.Payload), err // probably have DeserializeMessage return err
}

func (s *rtcm) NextFrame() (frame *Frame, err error) {
	for {
		select {
		case <-s.ctx.Done():
			return frame, s.ctx.Err()
		default:
		}

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
