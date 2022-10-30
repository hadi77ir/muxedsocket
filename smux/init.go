package smux

import "muxedsocket"

func init() {
	muxedsocket.GlobalCreators().ClientMuxers().Register("smux", ClientMuxer)
	muxedsocket.GlobalCreators().ServerMuxers().Register("smux", ServerMuxer)
}
