package yamux

import (
	"github.com/hadi77ir/muxedsocket"
	S "github.com/hashicorp/yamux"
)

// getConfig creates a new instance of Config, with prefilled values.
func getConfig(params muxedsocket.CommonParams) *S.Config {
	return &S.Config{
		EnableKeepAlive:    params.KeepalivePeriod <= 0 || params.MaxIdleTimeout <= params.KeepalivePeriod,
		KeepAliveInterval:  params.KeepalivePeriod,
		StreamCloseTimeout: params.MaxIdleTimeout,
	}
}
