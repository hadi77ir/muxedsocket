package tls

import (
	"bytes"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/utils"
	"strings"
)

const (
	ParamSNI                            = "sni"
	ParamNextProtos                     = "alpn"
	ParamHelloId                        = "profile"
	ParamCertificate                    = "cert"
	ParamPrivateKey                     = "key"
	ParamCertificatePin                 = "pin"
	ParamInsecure                       = "insecure"
	ParamClientCA                       = "clientca"
	CertificatePinDigestMethodSeparator = ":"
)

func LoadCertPoolFromParams(parameters utils.Parameters, paramName string) (*x509.CertPool, int, error) {
	pool := x509.NewCertPool()
	certificates, err := LoadCertsFromParams(parameters, paramName)
	if err != nil {
		return nil, 0, err
	}
	for _, cert := range certificates {
		pool.AddCert(cert)
	}
	return pool, len(certificates), nil
}

func LoadCertsFromParams(parameters utils.Parameters, paramName string) ([]*x509.Certificate, error) {
	paths, found := parameters.Get(paramName)
	pathsSplit := strings.Split(paths, muxedsocket.MultipleValuesSeparator)
	if found {
		certs := []*x509.Certificate{}
		for _, path := range pathsSplit {
			contents, err := utils.ReadFile(path)
			if err != nil {
				return nil, err
			}
			newCerts, err := x509.ParseCertificates(contents)
			if err != nil {
				return nil, err
			}
			certs = append(certs, newCerts...)
		}
		return certs, nil
	}
	return nil, nil
}

func LoadX509PairBytesFromParams(parameters utils.Parameters) (cert []byte, key []byte, err error) {
	keyPath, keyPathFound := parameters.Get(ParamPrivateKey)
	certPath, certPathFound := parameters.Get(ParamCertificate)
	if keyPathFound && !certPathFound {
		return nil, nil, muxedsocket.ErrMissingPart(ParamCertificate)
	}
	if !keyPathFound && certPathFound {
		return nil, nil, muxedsocket.ErrMissingPart(ParamPrivateKey)
	}
	if keyPathFound && certPathFound {
		return LoadX509PairBytes(certPath, keyPath)
	}
	// there is no error. as there was nothing to be loaded.
	return nil, nil, nil
}

func LoadX509PairBytes(certPath, keyPath string) (cert []byte, key []byte, err error) {
	key, err = utils.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}
	cert, err = utils.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}
	return
}
func LoadX509PairsBytesFromParams(parameters utils.Parameters) (certs [][]byte, keys [][]byte, err error) {
	keyPath, keyPathFound := parameters.Get(ParamPrivateKey)
	certPath, certPathFound := parameters.Get(ParamCertificate)
	if keyPathFound && !certPathFound {
		return nil, nil, muxedsocket.ErrMissingPart(ParamCertificate)
	}
	if !keyPathFound && certPathFound {
		return nil, nil, muxedsocket.ErrMissingPart(ParamPrivateKey)
	}
	keyPaths := strings.Split(keyPath, muxedsocket.MultipleValuesSeparator)
	certPaths := strings.Split(certPath, muxedsocket.MultipleValuesSeparator)
	if len(keyPaths) > len(certPaths) {
		return nil, nil, muxedsocket.ErrMissingPart(ParamCertificate)
	}
	if len(keyPaths) < len(certPaths) {
		return nil, nil, muxedsocket.ErrMissingPart(ParamPrivateKey)
	}
	pairCount := len(keyPaths)
	certs = make([][]byte, pairCount)
	keys = make([][]byte, pairCount)
	for i := 0; i < len(keyPaths); i++ {
		keys[i], certs[i], err = LoadX509PairBytes(certPaths[i], keyPaths[i])
		if err != nil {
			return nil, nil, err
		}
	}
	return
}

func GetCertificatePinningAndInsecure(parameters utils.Parameters) (vFunc PeerVerifierFunc, insecureBool bool, err error) {
	if insecure, found := parameters.Get(ParamInsecure); found {
		insecureBool, err = utils.ParseBool(insecure)
		if err != nil {
			return nil, false, err
		}
		vFunc = insecureVerifier
	}

	if pin, found := parameters.Get(ParamCertificatePin); found {
		verifyFunc, err := getPinVerificationFunc(pin)
		if err != nil {
			return nil, false, err
		}
		return verifyFunc, true, nil
	}
	return
}

type PeerVerifierFunc func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error

var certificateNotMatchingPinErr = errors.New("certificate fingerprint doesn't match with the pinned hash")

type CertificatePin struct {
	DigestFunc func([]byte) []byte
	Digest     []byte
}

func getPinVerificationFunc(pin string) (PeerVerifierFunc, error) {
	if pin != "" {
		pinsSplit := strings.Split(pin, muxedsocket.MultipleValuesSeparator)
		pins := make([]CertificatePin, len(pinsSplit))
		for i, pin := range pinsSplit {
			pinSplit := strings.SplitN(pin, CertificatePinDigestMethodSeparator, 2)
			pinBytes, err := hex.DecodeString(pinSplit[1])
			if err != nil {
				return nil, err
			}
			pins[i] = CertificatePin{DigestFunc: GetDigestFunc(pinSplit[0]), Digest: pinBytes}
			if pins[i].DigestFunc == nil {
				return nil, muxedsocket.ErrOpNotSupported
			}
		}
		return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if l := len(rawCerts); l != 1 {
				return fmt.Errorf("got len(rawCerts) = %d, wanted 1", l)
			}
			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return err
			}
			for _, pin := range pins {
				if bytes.Equal(pin.DigestFunc(cert.RawSubjectPublicKeyInfo), pin.Digest) {
					return nil
				}
			}
			return certificateNotMatchingPinErr
		}, nil
	}
	return nil, nil
}

func insecureVerifier(_ [][]byte, _ [][]*x509.Certificate) error {
	return nil
}

func GetSNIFromParams(parameters utils.Parameters) string {
	return utils.StringFromParameters(parameters, ParamSNI, "")
}
func GetNextProtosFromParams(parameters utils.Parameters) []string {
	return utils.MultiStringFromParameters(parameters, ParamNextProtos, nil)
}
