package updater

import (
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/g0dsCookie/ddnsd/protocols"
)

type Config struct {
	Name        string             `xml:"name,attr"`
	Protocol    protocols.Protocol `xml:"protocol,attr"`
	UpdateURL   string             `xml:"server"`
	Hosts       []string           `xml:"host"`
	Credentials CredentialsConfig  `xml:"credentials"`

	disabled    bool
	tempDisable time.Time
}

type CredentialsConfig struct {
	Username string `xml:"username,attr"`
	Password string `xml:"password,attr"`
}

func (c *Config) Apply() error {
	c.Name = strings.TrimFunc(c.Name, unicode.IsControl)
	if len(c.Name) == 0 {
		return errors.New("UpdateConfig name cannot be empty")
	}
	return nil
}

func (c *Config) Update() (bool, error) {
	if c.disabled {
		return true, nil
	}
	if c.tempDisable.After(time.Now()) {
		return false, nil
	}
	return protocols.Run(c.Protocol, c)
}

func (c *Config) GetUpdateURL() string { return c.UpdateURL }

func (c *Config) GetCredentials() (string, string) {
	return c.Credentials.Username, c.Credentials.Password
}

func (c *Config) GetHosts() []string {
	v := make([]string, len(c.Hosts))
	copy(v, c.Hosts)
	return v
}

func (c *Config) Disable() { c.disabled = true }

func (c *Config) TempDisable() { c.tempDisable = time.Now().Add(30 * time.Minute) }
