package utils

import (
	"encoding/base32"
	"encoding/base64"
	"os"
	"strings"
)

type Decoder interface {
	Decode(dst, src []byte) (n int, err error)
	DecodedLen(int) int
}

// ReadFile supports reading files at given directions, along base64 and base32 values.
func ReadFile(path string) ([]byte, error) {
	var decoder Decoder
	var toDecode []byte
	switch {
	case strings.HasPrefix(path, "base64:"):
		decoder = base64.StdEncoding
		toDecode = []byte(path[len("base64:"):])
		break
	case strings.HasPrefix(path, "base32:"):
		decoder = base32.StdEncoding
		toDecode = []byte(path[len("base32:"):])
		break
	}
	if decoder != nil {
		decoded := make([]byte, decoder.DecodedLen(len(path)))
		_, err := decoder.Decode(decoded, toDecode)
		if err != nil {
			return nil, err
		}
		return decoded, nil
	}
	return os.ReadFile(path)
}
