package nomux

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
)

type Stream struct {
	types.StreamConn
	dialer types.MuxStreamConnectFunc
}

func (s *Stream) CanRedial() bool {
	return s.dialer != nil
}

func (s *Stream) Redial() (types.Socket, error) {
	if s.dialer == nil {
		return nil, muxedsocket.ErrRedialNotSupported
	}
	return s.dialer()
}

func (s *Stream) StreamID() int {
	return 0
}

var _ types.MuxStream = &Stream{}

func WrapStream(conn types.StreamConn, dialFunc types.MuxStreamConnectFunc) types.MuxStream {
	return &Stream{StreamConn: conn, dialer: dialFunc}
}
