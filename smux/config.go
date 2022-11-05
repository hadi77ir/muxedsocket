package smux

import (
	"github.com/hadi77ir/muxedsocket"
	S "github.com/xtaci/smux"
)

const SupportedSMuxVersion = 2

// getConfig creates a new instance of Config, with prefilled values.
func getConfig(params muxedsocket.CommonParams) *S.Config {
	return &S.Config{
		Version:           SupportedSMuxVersion,
		KeepAliveDisabled: params.KeepalivePeriod <= 0 || params.MaxIdleTimeout <= params.KeepalivePeriod,
		KeepAliveInterval: params.KeepalivePeriod,
		KeepAliveTimeout:  params.MaxIdleTimeout,
	}
}
