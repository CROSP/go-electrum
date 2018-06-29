package main

import (
	"github.com/CROSP/go-electrum/electrum"
	"fmt"
)

func main() {
	peers, _ := electrum.FindElectrumPublicServers()
	fmt.Println(peers)
	var serverAddress *electrum.ServerAddress = nil
	var i = 0
	var node *electrum.Node
	for serverAddress == nil && i < len(*peers) {
		serverAddress = &(*peers)[i]
		node = electrum.NewNode()
		if err := node.ConnectTCP(serverAddress.GetAddressString()); err != nil {
			serverAddress = nil
		}
		i++
	}

	balance, _ := node.BlockchainAddressGetBalance("1FfmbHfnpaZjKFvyi1okTjJJusN455paPH")
	fmt.Println(balance)

}
