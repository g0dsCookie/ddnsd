package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/g0dsCookie/ddnsd/config"
	"github.com/g0dsCookie/ddnsd/log"
)

var (
	configFile string
	cron       bool
)

func init() {
	flag.StringVar(&configFile, "config", "ddnsd.xml", "the config file to use")
	flag.Parse()
}

func main() {
	if !flag.Parsed() {
		panic("flags were not parsed")
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}
	err = cfg.Apply()
	if err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	for {
		if len(cfg.Clients) == 0 {
			log.Warn("No active client left, closing...")
			return
		}

		select {
		case <-time.After(cfg.Global.Interval * time.Second):
			for i := 0; i < len(cfg.Clients); i++ {
				disabled, _ := cfg.Clients[i].Update()
				if disabled {
					cfg.Clients = append(cfg.Clients[:i], cfg.Clients[i+1:]...)
					i--
				}
			}
		case <-interrupt:
			log.Notice("Interrupt received, closing...")
			return
		}
	}
}
