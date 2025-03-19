//go:build utls
// +build utls

package tls

import (
	"errors"
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/utils"
	utls "github.com/refraction-networking/utls"
	"net"
	"strings"
)

var ErrProfileNotSupported = errors.New("profile not supported by uTLS library")

type TLSParams struct {
	Config        *utls.Config
	ClientHelloID utls.ClientHelloID
}

func ServerTLS(conn net.Conn, params any) net.Conn {
	config, _ := GetParamsUTLS(params)
	return utls.Server(conn, config)
}

func ClientTLS(conn net.Conn, params any) net.Conn {
	config, helloId := GetParamsUTLS(params)
	uconn := utls.UClient(conn, config, helloId)
	uconn.SetSNI(config.ServerName)
	return utils.WrapLazyHandshakingConn(uconn, uconn.HandshakeContext)
}

func ParseTLS(parameters utils.Parameters, isClient bool) (any, error) {
	tlsParams := &TLSParams{
		Config: &utls.Config{
			ServerName: GetSNIFromParams(parameters),
			NextProtos: GetNextProtosFromParams(parameters),
		},
	}

	if isClient {
		helloId, err := GetClientHelloIDFromParams(parameters)
		if err != nil {
			return nil, err
		}
		tlsParams.ClientHelloID = helloId

		verifierFunc, insecure, err := GetCertificatePinningAndInsecure(parameters)
		if err != nil {
			return nil, err
		}
		tlsParams.Config.InsecureSkipVerify = insecure
		tlsParams.Config.VerifyPeerCertificate = verifierFunc
	}

	if !isClient {
		clientCaPool, clientCaLen, err := LoadCertPoolFromParams(parameters, ParamClientCA)
		if err != nil {
			return nil, err
		}
		clientAuth := utls.NoClientCert
		if clientCaLen > 0 {
			clientAuth = utls.RequireAndVerifyClientCert
		}
		tlsParams.Config.ClientCAs = clientCaPool
		tlsParams.Config.ClientAuth = clientAuth
	}

	certs, err := LoadX509PairsFromParams(parameters)
	if err != nil {
		return nil, err
	}
	tlsParams.Config.Certificates = certs

	return tlsParams, nil
}

func LoadX509PairsFromParams(parameters utils.Parameters) ([]utls.Certificate, error) {
	certs, keys, err := LoadX509PairsBytesFromParams(parameters)
	if err != nil {
		return nil, err
	}
	pairs := make([]utls.Certificate, len(keys))
	for i := 0; i < len(keys); i++ {
		pairs[i], err = utls.X509KeyPair(certs[i], keys[i])
		if err != nil {
			return nil, err
		}
	}
	return pairs, nil
}

func LoadX509PairFromParams(parameters utils.Parameters) (utls.Certificate, error) {
	cert, key, err := LoadX509PairBytesFromParams(parameters)
	if err != nil {
		var zero utls.Certificate
		return zero, err
	}
	return utls.X509KeyPair(cert, key)
}

func GetClientHelloIDFromParams(parameters utils.Parameters) (utls.ClientHelloID, error) {
	profile, found := parameters.Get(ParamHelloId)
	if found {
		// let the user define client spec
		profileSplit := strings.Index(profile, muxedsocket.MultipleValuesSeparator)
		profileType := profile
		profileVer := ""
		if profileSplit != -1 {
			profileType = profile[:profileSplit]
			profileVer = profile[profileSplit+1:]
		}
		if profileVer == "" {
			profileType = strings.ToLower(profileType)
			switch profileType {
			case "chrome":
				return utls.HelloChrome_Auto, nil
			case "firefox":
				return utls.HelloFirefox_Auto, nil
			case "ios":
				return utls.HelloIOS_Auto, nil
			case "edge":
				return utls.HelloEdge_Auto, nil
			case "android":
				return utls.HelloAndroid_11_OkHttp, nil
			case "safari":
				return utls.HelloSafari_Auto, nil
			case "360":
				return utls.Hello360_Auto, nil
			case "qq":
				return utls.HelloQQ_Auto, nil
			default:
				return utls.ClientHelloID{}, ErrProfileNotSupported
			}
		}
		return utls.ClientHelloID{Client: profileType, Version: profileVer, Seed: nil}, nil
	}
	return utls.HelloGolang, nil
}

func GetParamsUTLS(params any) (*utls.Config, utls.ClientHelloID) {
	if castedParams, ok := params.(*TLSParams); ok && castedParams != nil {
		return castedParams.Config, castedParams.ClientHelloID
	}
	if castedParams, ok := params.(*utls.Config); ok {
		return castedParams, utls.HelloGolang
	}
	return nil, utls.HelloGolang
}
