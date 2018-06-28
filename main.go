package main

import (
	"github.com/CROSP/go-electrum/electrum"
	"fmt"
)

func main() {
	addres, _ := electrum.FindElectrumPublicServers()
	fmt.Println(addres[0])
}
