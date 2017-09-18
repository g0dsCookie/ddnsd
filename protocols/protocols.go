package protocols

import (
	"errors"

	"github.com/g0dsCookie/ddnsd/protocols/dyndns2"
)

type Protocol string

const (
	ProtocolDynDNS2 = "dyndns2"
)

type Config interface {
	GetUpdateURL() string
	GetHosts() []string
	GetCredentials() (string, string)
	Disable()
	TempDisable()
}

func Run(p Protocol, c Config) error {
	switch p {
	case ProtocolDynDNS2:
		return dyndns2.Update(c)

	default:
		return errors.New("unknown protocol " + string(p))
	}
}
