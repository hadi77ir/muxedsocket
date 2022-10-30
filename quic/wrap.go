package quic

import (
	Q "github.com/lucas-clemente/quic-go"
	"muxedsocket"
)

func wrapListener(listener Q.Listener) muxedsocket.MuxedListener {
	return &Listener{listener: listener}
}

func wrapConn(conn Q.Connection) (muxedsocket.DatagramCapableMuxedSocket, error) {
	return &Conn{conn: conn}, nil
}

func wrapStream(c Q.Connection, s Q.Stream) (muxedsocket.MuxStream, error) {
	return &Stream{
		stream:     s,
		localAddr:  muxedsocket.WrapAddr(c.LocalAddr(), int(s.StreamID())),
		remoteAddr: muxedsocket.WrapAddr(c.RemoteAddr(), int(s.StreamID())),
	}, nil
}
