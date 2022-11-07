package smux

import (
	"github.com/hadi77ir/muxedsocket"
	S "github.com/xtaci/smux"
)

const SupportedSMuxVersion = 2

// getConfig creates a new instance of Config, with prefilled values.
func getConfig(params muxedsocket.CommonParams) *S.Config {
	config := S.DefaultConfig()
	config.Version = SupportedSMuxVersion
	config.KeepAliveDisabled = params.KeepalivePeriod <= 0 || params.MaxIdleTimeout <= params.KeepalivePeriod
	config.KeepAliveInterval = params.KeepalivePeriod
	config.KeepAliveTimeout = params.MaxIdleTimeout
	return config
}
