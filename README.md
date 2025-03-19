# muxedsocket
This package enables use of multiple multiplexer and "stream-over-datagram" implementations with drop-in replacement
ability.

## Disclaimer / Project Status
This project is still under heavy development. So until further notice:
- Everything in this repository should be considered "Alpha-Quality Software" and therefore not used in production environments.
- API is unstable and is subject to change.
- Proposals for API and design changes that include drastic changes are welcome.

## Features
- Reusability. This library is built to be completely reusable.
- Extensible. You can extend functionality of the library and bring your own Multiplexer, Obfuscator, etc. 
  by registering them through `Creators`.
- Layered architecture. The library has multiple layers, which may be used independently or chained together.
- Compatible with interfaces defined in Go Standard Library (`net.Conn` and `net.PacketConn`). To adapt your code, there
  is little to zero modifications required.
- Simple URI Configuration. Two functions provide you the functionality with "Listen" and "Dial" functionality,
  and they just take a URI.
- Stream-over-Packets: Supports Packet-oriented connections for networks where UDP is not monitored.
- Packets-over-streams: Supports "Packets-over-streams" when Packet-oriented connections are not available and 
  Stream-oriented connections are unstable or limited.

## Samples
For samples, take a look at the tests and the `everest` project.

## License
Apache License 2.0 - See [LICENSE](LICENSE)
