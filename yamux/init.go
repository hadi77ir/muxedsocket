package yamux

import "muxedsocket"

func init() {
	muxedsocket.GlobalCreators().ClientMuxers().Register("yamux", ClientMuxer)
	muxedsocket.GlobalCreators().ServerMuxers().Register("yamux", ServerMuxer)
}
