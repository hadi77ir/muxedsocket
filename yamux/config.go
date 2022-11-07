package yamux

import (
	"github.com/hadi77ir/muxedsocket"
	S "github.com/hashicorp/yamux"
)

// getConfig creates a new instance of Config, with prefilled values.
func getConfig(params muxedsocket.CommonParams) *S.Config {
	config := S.DefaultConfig()
	config.ConnectionWriteTimeout = params.MaxIdleTimeout
	config.EnableKeepAlive = params.KeepalivePeriod <= 0 || params.MaxIdleTimeout <= params.KeepalivePeriod
	config.KeepAliveInterval = params.KeepalivePeriod
	config.StreamCloseTimeout = params.MaxIdleTimeout
	return config
}
