# Multi-route Packet-oriented Tunnel

Recently, our adversaries behind the Great FireWall implemented methods for detecting "WebSocket" and HTTP/2 connections
which are being actively used for proxying and tunneling: They watch long-living TLS connections that are having too
much data transfer (both receive and send) closely and increase the server's probability score of tunneling. This
proposal, provides a way that (at least in theory) may make it difficult for them to detect tunneling.

## Streams-over-Packets: ARQ connection on top of Packet transports

A complete ARQ implementation like QUIC may help us use UDP and other protocols (such as ICMP and DNS) for streaming
connections, but this is still distinctly recognizable traffic. How we may hide our traffic as innocuous
WebSocket-over-TLS connections?

## Packets-over-Streams: Short-lived TCP connections as Packet transports

We may turn WebSocket or any other TCP-based protocol (like SSH) into a packet transport. While it is not too much
efficient, it will work. We may send packets in a specific format to the server, where it can decode packets and turn
them into UDP packets and send packets to UDP server and such for replies from UDP server.

But as told in the previous section, the goal is to have a streaming, stable connection like TCP. This is where QUIC
comes in handy. We transport QUIC inside TCP. But a single connection is not that much efficient. So we may add more
concurrent connections, which may make it more UDPish: unstable, unordered, different paths with different performances,
etc.

## One step further!

These are all performed between one client and one specific server, which will be easily recognized by the firewall and
the censor party and then will be blocked. We may add one step of uncertainty to this: Add packet-relay servers to the
chain.

By adding packet-relay servers, streaming traffic will go through multiple different packet-transport connections with
different destinations "that all relay traffic to one final destination".