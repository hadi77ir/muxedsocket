package yamux

import (
	"github.com/hadi77ir/muxedsocket"
	Y "github.com/hashicorp/yamux"
	"net"
	"time"
)

var _ muxedsocket.MuxStream = &Stream{}

type Stream struct {
	stream     *Y.Stream
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (s *Stream) Read(b []byte) (n int, err error) {
	return s.stream.Read(b)
}

func (s *Stream) Write(b []byte) (n int, err error) {
	return s.stream.Write(b)
}

func (s *Stream) Close() error {
	return s.stream.Close()
}

func (s *Stream) LocalAddr() net.Addr {
	return s.localAddr
}

func (s *Stream) RemoteAddr() net.Addr {
	return s.remoteAddr
}

func (s *Stream) SetDeadline(t time.Time) error {
	return s.SetDeadline(t)
}

func (s *Stream) SetReadDeadline(t time.Time) error {
	return s.SetReadDeadline(t)
}

func (s *Stream) SetWriteDeadline(t time.Time) error {
	return s.SetWriteDeadline(t)
}

func (s *Stream) StreamID() int {
	return int(s.stream.StreamID())
}
