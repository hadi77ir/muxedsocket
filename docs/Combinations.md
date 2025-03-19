## Available Parts

- SHP: Traffic Shaping for Streaming Connections: HTTP/2, WebSocket, gRPC,...
- CON: Concrete Streaming Connections. Simply TCP, Unix socket or Pipe
- PKT: Packet Connections, like: UDP, ICMP, etc.
- POS: Packets over Streams: PoS, SPoS (Session-based PoS), PSPoS (Parallel SPoS)
- MUX: Stream Multiplexer: smux, yamux
- SOP: Stream over Packets: KCP
- OBF: Stream obfuscators: TLS, uTLS, ...
- POB: Packet obfuscators: (are there any out there?)
- PAIO: All-in-one solutions for Packet-based connections (Multiplexer + Obfuscator + Stream over Packets): QUIC
- SAIO: All-in-one solutions for Stream-based connections (Multiplexer + Traffic Shaper): HTTP/2 Cleartext

## Possible Combinations

### Streaming Transports (ST)

- (OBF+)CON 
- (OBF+)SOP+PT

### Packet Transports (PT)

- PKT
- POB+PKT
- POS+ST
- POB+POS+ST

### Muxed Transport (MT)

- (OBF+)MUX+ST
- (OBF+)SAIO+ST
- (OBF+)PAIO+PT
