package muxedsocket

import (
	"crypto/tls"
	"net"
	"net/url"
	"strings"
)

func CreateListenerWithRegistry(creators *Creators, scheme string) (AddrMuxListener, error) {
	if listener, found := creators.AddrMuxListeners().Get(scheme); found {
		return listener, nil
	}
	schemeParts := strings.Split(scheme, "+")

	// try finding the mux listener in the scheme
	schemeParts, _, muxListener, muxListenerFound := findAndRemove(schemeParts, creators.MuxListeners())

	// try finding the packet transport in the scheme
	schemeParts, channelDialerKey, channelDialer, channelDialerFound := findAndRemove(schemeParts, creators.ChannelDialers())
	if !channelDialerFound {
		channelDialer, channelDialerFound = creators.ChannelDialers().Get(defaultPacketTransport)
	}
	if muxListenerFound {
		// there should not be anything remaining
		if len(schemeParts) > 0 {
			return nil, SchemeNotSupported
		}
		return func(addr string, config *tls.Config) (MuxedListener, error) {
			parsed, err := url.Parse(addr)
			if err != nil {
				return nil, err
			}
			transport, err := channelDialer(parsed.Host)
			if err != nil {
				return nil, err
			}
			params, err := GetServerParamsFromURL(channelDialerKey, parsed)
			if err != nil {
				return nil, err
			}
			return muxListener(transport, config, params)
		}, nil
	}

	// try finding the client muxer in the scheme
	schemeParts, _, serverMuxer, serverMuxerFound := findAndRemove(schemeParts, creators.ServerMuxers())
	if !serverMuxerFound {
		serverMuxer, serverMuxerFound = creators.ServerMuxers().Get(defaultMuxer)
	}

	// try finding the stream adapter in the scheme
	schemeParts, streamAdapterKey, streamAdapter, streamAdapterFound := findAndRemove(schemeParts, creators.ServerStreamAdapters())

	// check if it has "secure" or "tls"
	schemeParts, secureConn := findAndRemoveTLS(schemeParts)

	// no stream adapter, so use TCP
	if !streamAdapterFound {
		if len(schemeParts) > 1 || (len(schemeParts) == 1 && schemeParts[0] != defaultStreamTransport) {
			return nil, SchemeNotSupported
		}
		return func(addr string, config *tls.Config) (MuxedListener, error) {
			parsed, err := url.Parse(addr)
			if err != nil {
				return nil, err
			}
			params, err := GetServerParamsFromURL(defaultStreamTransport, parsed)
			if err != nil {
				return nil, err
			}
			var conn net.Listener
			if !secureConn {
				conn, err = net.Listen(defaultStreamTransport, parsed.Host)
				if err != nil {
					return nil, err
				}
			} else {
				conn, err = tls.Listen(defaultStreamTransport, parsed.Host, config)
				if err != nil {
					return nil, err
				}
			}
			muxed, err := serverMuxer(conn, params)
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

	return func(addr string, config *tls.Config) (MuxedListener, error) {
		parsed, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		params, err := GetServerParamsFromURL(streamAdapterKey, parsed)
		if err != nil {
			return nil, err
		}
		transport, err := channelDialer(parsed.Host)
		if err != nil {
			return nil, err
		}
		conn, err := streamAdapter(transport)
		if err != nil {
			_ = transport.Close()
			return nil, err
		}
		if secureConn {
			conn = tls.NewListener(conn, config)
		}
		muxed, err := serverMuxer(conn, params)
		if err != nil {
			_ = conn.Close()
			_ = transport.Close()
			return nil, err
		}
		return muxed, nil
	}, nil
}

func CreateListener(scheme string) (AddrMuxListener, error) {
	return CreateListenerWithRegistry(creators, scheme)
}

func ListenURI(uri string, config *tls.Config) (MuxedListener, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	listener, err := CreateListener(parsed.Scheme)
	if err != nil {
		return nil, err
	}
	return listener(uri, config)
}

func ListenURIWithRegistry(creators *Creators, uri string, config *tls.Config) (MuxedListener, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	listener, err := CreateListenerWithRegistry(creators, parsed.Scheme)
	if err != nil {
		return nil, err
	}
	return listener(uri, config)
}

func GetServerParamsFromURL(_ string, addr *url.URL) (*ServerParams, error) {
	return &ServerParams{CommonParams: GetCommonParamsFromURL(addr)}, nil
}
