package quic

import "muxedsocket"

func init() {
	muxedsocket.GlobalCreators().MuxDialers().Register("quic", Dial)
	muxedsocket.GlobalCreators().MuxListeners().Register("quic", Listen)
}
