package main

import (
	"fmt"

	"github.com/g0dsCookie/ddnsd/config"
	"github.com/g0dsCookie/ddnsd/log"
)

func main() {
	cfg, err := config.LoadConfig("ddnsd.xml")
	if err != nil {
		panic(err)
	}
	err = cfg.Apply()
	if err != nil {
		panic(err)
	}

	for _, v := range cfg.Clients {
		if err := v.Update(); err != nil {
			log.Err(fmt.Sprintf("%v failed: %v", v.Name, err.Error()))
		}
		log.Notice(fmt.Sprintf("%v updated", v.Name))
	}
}
