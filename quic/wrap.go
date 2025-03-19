package quic

import (
	"github.com/hadi77ir/muxedsocket/types"
	Q "github.com/lucas-clemente/quic-go"
)

func wrapListener(listener Q.Listener) types.MuxedListener {
	return &Listener{listener: listener}
}

func wrapConn(conn Q.Connection) (types.DatagramCapableMuxedSocket, error) {
	return &Conn{conn: conn}, nil
}

func wrapStream(c Q.Connection, s Q.Stream) (types.MuxStream, error) {
	return &Stream{
		stream:     s,
		localAddr:  types.WrapAddr(c.LocalAddr(), int(s.StreamID())),
		remoteAddr: types.WrapAddr(c.RemoteAddr(), int(s.StreamID())),
	}, nil
}
