# muxedsocket
This package enables use of multiple multiplexer and "stream-over-datagram" implementations with drop-in replacement
ability.

## Supports...?
Supported implementations of multiplexer are:
- SPDY
- smux
- yamux

Supported stream connection implementations are:
- Stream-oriented protocols:
  - TCP
- Stream-over-packets:
  - KCP
  - KCP Secure (KCP with TLS)

Complete solutions:
- QUIC

When using anything other than TCP, you may use a custom `PacketConn` implementation such as `ICMPChannel` to transport 
traffic over some exotic channel of your choice.

## Design Choices
There are a number of design choices that I rather explain here.

### A note on contexts
I would like to use `Context`s but the inconsistencies led to the decision to drop them in all signatures.
As we already have `Close` for signaling the end of listening or connection lifetime, I think there is no real
compensation for trying to implement them.

## License
Apache License 2.0 - See [LICENSE](LICENSE)
