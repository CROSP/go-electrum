package main

import (
	"log"

	"github.com/CROSP/go-electrum/irc"
)

func main() {
	log.Println(irc.FindElectrumServers())
}
