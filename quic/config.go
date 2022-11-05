package quic

import (
	"github.com/hadi77ir/muxedsocket"
	Q "github.com/lucas-clemente/quic-go"
)

// getConfig creates a new instance of Config, with prefilled values.
func getConfig(params muxedsocket.CommonParams) *Q.Config {
	return &Q.Config{
		// datagrams are supported
		EnableDatagrams: true,
		// only allow Version 2
		Versions: []Q.VersionNumber{Q.Version2},
		// the default 100 is a bit low, but 1000 seems a good choice
		MaxIncomingStreams: 1000,
		MaxIdleTimeout:     params.MaxIdleTimeout,
		KeepAlivePeriod:    params.KeepalivePeriod,
	}
}
