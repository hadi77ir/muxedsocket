package muxedsocket

import (
	"strings"
	"time"
)

func GetSchemeParts(scheme string) []string {
	return strings.Split(scheme, "+")
}

const (
	// ParamDialTimeout defines how much it takes before dialing times out.
	ParamDialTimeout = "timeout"
	// ParamKeepAlive defines if keep-alive messages have to be sent and at what interval. May contain "false" or duration.
	ParamKeepAlive = "keepalive"
	// ParamDPD means time that has to pass after a keep-alive has been sent and no response was received to assume the
	// connection is dead. (dead peer detection)
	ParamDPD = "dpd"
)

const (
	DefaultDialTimeout = time.Duration(5) * time.Second
	DefaultKeepAlive   = time.Duration(15) * time.Second
)

// MultipleValuesSeparator is used to split a parameter value into multiple values.
const MultipleValuesSeparator = ","
