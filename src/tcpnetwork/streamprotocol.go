package tcpnetwork

import (
	"encoding/binary"
	"errors"
)

const (
	KStreamProtocol4HeaderLength = 4
	KStreamProtocol2HeaderLength = 2
	kStreamProtocolBufferCap     = 1024 * 64
)

var (
	kStreamProtocolErrFail        = errors.New("fail")
	kStreamProtocolErrPackNotFull = errors.New("package not fulled.")
	kStreamProtocolErrSerWrong    = errors.New("serialize occur wrong.")
)

func getStreamMaxLength(headerBytes uint32) uint64 {
	return 1<<(8*headerBytes) - 1
}

type IStreamProtocol interface {
	Serialize(body interface{}) ([]byte, error)
	Unserialize(bin []byte) (interface{}, error)
}

// StreamProtocol4
// Binary format : | 4 byte (total stream length) | data ... (total stream length - 4) |
//	implement default stream protocol
//	stream protocol interface for 4 bytes header
type StreamProtocol struct {
	Buffer       [kStreamProtocolBufferCap]byte
	Length       int
	headerLength int
	handler      *IStreamProtocol
}

func NewStreamProtocol(headerLength int) *StreamProtocol {
	s := &StreamProtocol{
		Length:       0,
		headerLength: headerLength,
	}
	return s
}

func (s *StreamProtocol) SetHandler(handler *IStreamProtocol) {
	s.handler = handler
}

func (s *StreamProtocol) Unserialize() (int, interface{}, error) {
	var n int
	var o interface{}
	var err error

	if s.Length > s.headerLength {
		if s.headerLength == KStreamProtocol2HeaderLength {
			length := binary.BigEndian.Uint16(s.Buffer[0:])
			if int(length)+KStreamProtocol2HeaderLength <= s.Length {
				o, err = s.handler.Unserialize(s.Buffer[2:])
				return length, o, err
			} else {
				o = nil
				err = kStreamProtocolErrFail
				return length, o, err
			}
		}
	} else {
		n = 0
		o = nil
		err = kStreamProtocolErrFail
		return n, o, err
	}
}

func (s *StreamProtocol) Serialize(body interface{}) ([]byte, error) {
	if nil == body {
		return nil, kStreamProtocolErrFail
	}
	buf, err := s.handler.Serialize(body)
	if err == nil {
		buffer := make([]byte, s.headerLength+len(buf))
		if s.headerLength == KStreamProtocol2HeaderLength {
			binary.BigEndian.PutUint16(buffer[0:], len(buf))
			copy(buffer[2:], buf)
		} else {
			binary.BigEndian.PutUint32(buffer[0:], len(buf))
			copy(buffer[4:], buf)
		}
		return buffer
	}
	return nil, kStreamProtocolSerWrong
}
