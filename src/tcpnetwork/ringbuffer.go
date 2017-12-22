package tcpnetwork

import (
	"encoding/binary"
	"errors"
)

const (
	kRingBufferBufferCap = 1024 * 64
)

var (
	kRingBufferErrFail        = errors.New("fail")
	kRingBufferErrPackNotFull = errors.New("package not fulled.")
	kRingBufferErrSerWrong    = errors.New("serialize occur wrong.")
)

type RingBuffer struct {
	conn   Connection
	buffer []byte
	head   int
	tail   int
	size   int
}

func NewRingBuffer(conn Connection) *RingBuffer {
	s := &RingBuffer{
		conn:   conn,
		buffer: make([]byte, kRingBufferBufferCap*2),
		head:   0,
		tail:   0,
		size:   kRingBufferBufferCap,
	}
	return s
}

func (s *RingBuffer) ext() {
	size := s.size * 2 * 2
	buffer := make([]byte, size)
	if s.head >= s.tail {
		copy(buffer[0:], s.buffer[0:s.size])
	} else {
		used := s.bytes_used()
		copy(buffer[0:], s.buffer[:s.head])
		copy(buffer[s.head:], s.buffer[s.tail:s.size])
		s.head = 0
		s.tail = used
		s.size = size
	}
	s.buffer = buffer
}

func (s *RingBuffer) end() int {
	return s.size
}

func (s *RingBuffer) cap() int {
	return s.size - 1
}

func (s *RingBuffer) freeBytes() int {
	if s.head >= s.tail {
		return s.cap() - (s.head - s.tail)
	} else {
		return s.tail - s.head - 1
	}
}

func (s *RingBuffer) usedBytes() int {
	return s.cap() - s.bytes_free()
}

func (s *RingBuffer) isFull() bool {
	if s.freeBytes() == 0 {
		return true
	} else {
		return false
	}
}

func (s *RingBuffer) isEmpty() {
	if s.freeBytes() == s.cap() {
		return true
	} else {
		return false
	}
}

func (s *RingBuffer) correct() {
	if s.isEmpty() {
		return
	}
	if s.head >= s.tail {
		copy(s.buffer[0:], s.buffer[s.tail:s.head])
	} else {
		copy(s.buffer[s.end():], s.buffer[0:s.head])
		copy(s.buffer[s.buffer[0:], s.buffer[s.tail:s.end()])
		copy(s.buffer[s.end()-s.tail:], s.buffer[s.end():s.end()+s.head])
		head := s.usedBytes()
		s.head = head
		s.tail = 0
	}
}

func (s *RingBuffer) Read(p []byte) (int, error) {
	var n int
	var err error
	for {
		if s.isFull() {
			s.ext()
		}
		if s.head >= s.tail {
			n, err = s.conn.conn.Read(s.buffer[s.head:s.size])
			if err == nil {
				s.head += n
			} else {
				break
			}
			if s.head == s.size {
				s.head = 0
			}
		} else {

		}
	}
	if s.usedBytes() >= 2 {
		if s.head >= s.tail {
			len := binary.BigEndian.Uint16(s.buffer[s.tail])
			s.tail += 2
			if s.usedBytes() >= len {
				p = s.buffer[s.tail:s.tail+len]
				s.tail += len
				n = len
				return n
			} else {
				n = 0
				return 0, kRingBufferErrPackNotFull
			}
		} else {
			if s.end() - s.tail >= 2 {
				len := binary.BigEndian.Uint16(s.buffer[s.tail])
				s.tail += 2
				
			}
		}
	} else {
		return 0, kRingBufferErrPackNotFull
	}
}

func (s *RingBuffer) Write(p []byte) (int, error) {
	if nil == body {
		return kRingBufferErrFail
	}
	if s.bytes_free() < 2 {
		s.ext()
	}
	if s.head >= s.tail {
		if s.end()-s.head >= 2 {
			binary.BigEndian.PutUint16(s.buffer[s.head:], len(body))
			s.head += 2
		} else {

		}
	}

	buf, err := s.handler.Serialize(body)
	if err == nil {
		buffer := make([]byte, s.headerLength+len(buf))
		if s.headerLength == KRingBuffer2HeaderLength {
			binary.BigEndian.PutUint16(buffer[0:], len(buf))
			copy(buffer[2:], buf)
		} else {
			binary.BigEndian.PutUint32(buffer[0:], len(buf))
			copy(buffer[4:], buf)
		}
		return buffer
	}
	_, err = c.conn.Write(b)

	return nil, kRingBufferSerWrong
}
