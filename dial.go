package muxedsocket

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var SchemeNotSupported = errors.New("scheme not supported")

const defaultPacketTransport = "udp"
const defaultStreamTransport = "tcp"
const defaultMuxer = "yamux"

// remove an element at given index. don't use if order matters in slice.
// from https://stackoverflow.com/a/37335777
func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findAndRemove[T any](keys []string, registry *Registry[T]) ([]string, string, T, bool) {
	for i := 0; i < len(keys); i++ {
		if value, found := registry.Get(keys[i]); found {
			keys = remove(keys, i)
			return keys, keys[i], value, true
		}
	}
	var none T
	return keys, "", none, false
}
func findAndRemoveTLS(parts []string) ([]string, bool) {
	for i := 0; i < len(parts); i++ {
		if parts[i] == "secure" || parts[i] == "tls" {
			return remove(parts, i), true
		}
	}
	return parts, false
}

// CreateDialer enables to create a dialer function for the provided scheme and cache it. If you dial a scheme frequently,
// then you probably should use it and cache its result.
func CreateDialer(scheme string) (AddrMuxDialer, error) {
	return CreateDialerWithRegistry(creators, scheme)
}

// CreateDialerWithRegistry is the same as CreateDialer but with 2 arguments: it takes your registry of creators instead of default.
func CreateDialerWithRegistry(creators *Creators, scheme string) (AddrMuxDialer, error) {
	if dialer, found := creators.AddrMuxDialers().Get(scheme); found {
		return dialer, nil
	}
	schemeParts := strings.Split(scheme, "+")

	// try finding the mux dialer in the scheme
	schemeParts, _, muxDialer, muxDialerFound := findAndRemove(schemeParts, creators.MuxDialers())

	// try finding the packet transport in the scheme
	schemeParts, channelDialerKey, channelDialer, channelDialerFound := findAndRemove(schemeParts, creators.ChannelDialers())
	if !channelDialerFound {
		channelDialer, channelDialerFound = creators.ChannelDialers().Get(defaultPacketTransport)
	}
	if muxDialerFound {
		// there should not be anything remaining
		if len(schemeParts) > 0 {
			return nil, SchemeNotSupported
		}
		return func(addr string, config *tls.Config) (MuxedSocket, error) {
			parsed, err := url.Parse(addr)
			if err != nil {
				return nil, err
			}
			transport, err := channelDialer(parsed.Host)
			if err != nil {
				return nil, err
			}
			params, err := GetClientParamsFromURL(channelDialerKey, parsed)
			if err != nil {
				return nil, err
			}
			return muxDialer(transport, config, params)
		}, nil
	}

	// try finding the client muxer in the scheme
	schemeParts, _, clientMuxer, clientMuxerFound := findAndRemove(schemeParts, creators.ClientMuxers())
	if !clientMuxerFound {
		clientMuxer, clientMuxerFound = creators.ClientMuxers().Get(defaultMuxer)
	}

	// try finding the stream adapter in the scheme
	schemeParts, streamAdapterKey, streamAdapter, streamAdapterFound := findAndRemove(schemeParts, creators.ClientStreamAdapters())

	// check if it has "secure" or "tls"
	schemeParts, secureConn := findAndRemoveTLS(schemeParts)

	// no stream adapter, so use TCP
	if !streamAdapterFound {
		if len(schemeParts) > 1 || (len(schemeParts) == 1 && schemeParts[0] != defaultStreamTransport) {
			return nil, SchemeNotSupported
		}
		return func(addr string, config *tls.Config) (MuxedSocket, error) {
			parsed, err := url.Parse(addr)
			if err != nil {
				return nil, err
			}
			params, err := GetClientParamsFromURL(defaultStreamTransport, parsed)
			if err != nil {
				return nil, err
			}
			var conn net.Conn
			if !secureConn {
				conn, err = net.Dial(defaultStreamTransport, parsed.Host)
				if err != nil {
					return nil, err
				}
			} else {
				conn, err = tls.Dial(defaultStreamTransport, parsed.Host, config)
				if err != nil {
					return nil, err
				}
			}
			muxed, err := clientMuxer(conn, params)
			if err != nil {
				_ = conn.Close()
				return nil, err
			}
			return muxed, nil
		}, nil
	}

	if len(schemeParts) > 0 {
		return nil, SchemeNotSupported
	}

	return func(addr string, config *tls.Config) (MuxedSocket, error) {
		parsed, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		params, err := GetClientParamsFromURL(streamAdapterKey, parsed)
		if err != nil {
			return nil, err
		}
		transport, err := channelDialer(parsed.Host)
		if err != nil {
			return nil, err
		}
		conn, err := streamAdapter(transport, params.RemoteAddr.String())
		if err != nil {
			_ = transport.Close()
			return nil, err
		}
		if secureConn {
			conn = tls.Client(conn, config)
		}
		muxed, err := clientMuxer(conn, params)
		if err != nil {
			_ = conn.Close()
			_ = transport.Close()
			return nil, err
		}
		return muxed, nil
	}, nil
}

func GetClientParamsFromURL(transportKey string, addr *url.URL) (*ClientParams, error) {
	keepalive, _ := time.ParseDuration(addr.Query().Get("keepalive"))
	timeout, _ := time.ParseDuration(addr.Query().Get("timeout"))
	remoteAddr, err := GetAddrByTransportType(transportKey, addr.Host)
	if err != nil {
		return nil, err
	}
	return &ClientParams{
		CommonParams: CommonParams{
			KeepalivePeriod: keepalive,
			MaxIdleTimeout:  timeout,
		},
		RemoteAddr: remoteAddr,
	}, nil
}

func GetAddrByTransportType(transportKey string, addr string) (NetAddrPort, error) {
	sep := strings.Index(addr, ":")
	var ipStr string
	var portStr string
	if sep != -1 {
		ipStr = addr[:sep]
		portStr = addr[sep+1:]
	} else {
		ipStr = addr
		portStr = "0"
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, net.InvalidAddrError("invalid ip address")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	if transportKey == "udp" {
		return &net.UDPAddr{IP: ip, Port: port}, nil
	}
	if transportKey == "tcp" {
		return &net.TCPAddr{IP: ip, Port: port}, nil
	}

	return &netAddr{
		host:    addr,
		port:    uint16(port),
		network: transportKey,
	}, nil
}

// DialURI provides a standard way of dialing an address.
func DialURI(uri string, config *tls.Config) (MuxedSocket, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	dialer, err := CreateDialer(parsed.Scheme)
	if err != nil {
		return nil, err
	}
	return dialer(uri, config)
}

// DialURIWithRegistry is the same as DialURI except it supports providing custom registry of creators.
func DialURIWithRegistry(creators *Creators, uri string, config *tls.Config) (MuxedSocket, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	dialer, err := CreateDialerWithRegistry(creators, parsed.Scheme)
	if err != nil {
		return nil, err
	}
	return dialer(uri, config)
}
