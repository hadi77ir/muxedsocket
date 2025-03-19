package stream

import (
	"github.com/hadi77ir/muxedsocket/types"
	"net"
)

func NewUnixSocketImplementation() types.StreamConnImplementation {
	return WrapStandardImplementation("unix", net.DialTimeout, net.Listen)
}
