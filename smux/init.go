package smux

import "github.com/hadi77ir/muxedsocket"

func init() {
	muxedsocket.GlobalCreators().ClientMuxers().Register("smux", ClientMuxer)
	muxedsocket.GlobalCreators().ServerMuxers().Register("smux", ServerMuxer)
}
