package muxedsocket

import (
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
)

// Creators structure will contain registries for different implementation factories
type Creators struct {
	packetConns       *utils.Registry[types.PacketConnImplementation]
	streamConns       *utils.Registry[types.StreamConnImplementation]
	streamAdapters    *utils.Registry[types.StreamAdapterImplementation]
	packetAdapters    *utils.Registry[types.PacketAdapterImplementation]
	streamObfuscators *utils.Registry[types.StreamObfuscatorImplementation]
	packetObfuscators *utils.Registry[types.PacketObfuscatorImplementation]
	packetSolutions   *utils.Registry[types.PacketSolutionImplementation]
	streamSolutions   *utils.Registry[types.StreamSolutionImplementation]
	addrSolutions     *utils.Registry[types.AddrSolutionImplementation]
}

func (c *Creators) PacketConns() *utils.Registry[types.PacketConnImplementation] {
	return c.packetConns
}

func (c *Creators) StreamConns() *utils.Registry[types.StreamConnImplementation] {
	return c.streamConns
}

func (c *Creators) StreamAdapters() *utils.Registry[types.StreamAdapterImplementation] {
	return c.streamAdapters
}

func (c *Creators) PacketAdapters() *utils.Registry[types.PacketAdapterImplementation] {
	return c.packetAdapters
}

func (c *Creators) StreamObfuscators() *utils.Registry[types.StreamObfuscatorImplementation] {
	return c.streamObfuscators
}

func (c *Creators) PacketObfuscators() *utils.Registry[types.PacketObfuscatorImplementation] {
	return c.packetObfuscators
}

func (c *Creators) PacketSolutions() *utils.Registry[types.PacketSolutionImplementation] {
	return c.packetSolutions
}

func (c *Creators) StreamSolutions() *utils.Registry[types.StreamSolutionImplementation] {
	return c.streamSolutions
}

func (c *Creators) AddrSolutions() *utils.Registry[types.AddrSolutionImplementation] {
	return c.addrSolutions
}

var creators = NewCreators()

func GlobalCreators() *Creators {
	return creators
}

func NewCreators() *Creators {
	return &Creators{
		packetConns:       &utils.Registry[types.PacketConnImplementation]{},
		streamConns:       &utils.Registry[types.StreamConnImplementation]{},
		streamAdapters:    &utils.Registry[types.StreamAdapterImplementation]{},
		packetAdapters:    &utils.Registry[types.PacketAdapterImplementation]{},
		streamObfuscators: &utils.Registry[types.StreamObfuscatorImplementation]{},
		packetObfuscators: &utils.Registry[types.PacketObfuscatorImplementation]{},
		packetSolutions:   &utils.Registry[types.PacketSolutionImplementation]{},
		streamSolutions:   &utils.Registry[types.StreamSolutionImplementation]{},
		addrSolutions:     &utils.Registry[types.AddrSolutionImplementation]{},
	}
}
