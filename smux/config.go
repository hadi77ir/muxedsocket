package smux

import (
	S "github.com/xtaci/smux"
	"muxedsocket"
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
