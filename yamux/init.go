package yamux

import "github.com/hadi77ir/muxedsocket"

func init() {
	muxedsocket.GlobalCreators().ClientMuxers().Register("yamux", ClientMuxer)
	muxedsocket.GlobalCreators().ServerMuxers().Register("yamux", ServerMuxer)
}
