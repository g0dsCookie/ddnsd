package config

import (
	"encoding/xml"
	"io/ioutil"
	"time"

	"github.com/g0dsCookie/ddnsd/log"
	"github.com/g0dsCookie/ddnsd/updater"
)

type Config struct {
	XMLName xml.Name          `xml:"ddnsd"`
	Global  GlobalConfig      `xml:"global"`
	Clients []*updater.Config `xml:"client"`
}

type GlobalConfig struct {
	Interval time.Duration `xml:"interval"`
	Log      log.LogConfig `xml:"log"`
}

func LoadConfig(file string) (c Config, err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = xml.Unmarshal(b, &c)
	return
}

func (c Config) Apply() error {
	if err := c.Global.Apply(); err != nil {
		return err
	}
	for _, v := range c.Clients {
		if err := v.Apply(); err != nil {
			return err
		}
	}
	return nil
}

func (c GlobalConfig) Apply() error {
	return c.Log.Apply()
}
