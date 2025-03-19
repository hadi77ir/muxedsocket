package stream

import (
	"github.com/hadi77ir/muxedsocket/types"
	"net"
)

func NewTCPImplementation() types.StreamConnImplementation {
	impl := WrapStandardImplementation("tcp", net.DialTimeout, net.Listen)
	AddKeepAliveHook(impl)
	return impl
}
