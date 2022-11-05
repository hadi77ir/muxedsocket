package quic

import "github.com/hadi77ir/muxedsocket"

func init() {
	muxedsocket.GlobalCreators().MuxDialers().Register("quic", Dial)
	muxedsocket.GlobalCreators().MuxListeners().Register("quic", Listen)
}
