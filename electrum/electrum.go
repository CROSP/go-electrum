// Package electrum is a electrum bitcoin client.
package electrum

import (
	"fmt"
	"math/rand"
	"time"
	"net"
	"os"
	"log"
	"encoding/json"
)

const (
	ConnectionTimeout = 5 * time.Second
)
type ServerAddress struct {
	hostname string
	portT int
	portS int
}
var DefaultServers = []ServerAddress{
	{hostname: "erbium1.sytes.net", portS: 50001, portT: 50002},
	{hostname: "ecdsa.net", portS: 50001, portT: 110},
	{hostname: "gh05.geekhosters.com", portS: 50001, portT: 50002},
	{hostname: "eVPS.hsmiths.com", portS: 50001, portT: 50002},
	{hostname: "electrum.anduck.net", portS: 50001, portT: 50002},
	{hostname: "electrum.no-ip.org", portS: 50001, portT: 50002},
	{hostname: "electrum.be", portS: 50001, portT: 50002},
	{hostname: "helicarrier.bauerj.eu", portS: 50001, portT: 50002},
	{hostname: "elex01.blackpole.online", portS: 50001, portT: 50002},
	{hostname: "electrumx.not.fyi", portS: 50001, portT: 50002},
	{hostname: "node.xbt.eu", portS: 50001, portT: 50002},
	{hostname: "kirsche.emzy.de", portS: 50001, portT: 50002},
	{hostname: "electrum.villocq.com", portS: 50001, portT: 50002},
	{hostname: "us11.einfachmalnettsein.de", portS: 50001, portT: 50002},
	{hostname: "electrum.trouth.net", portS: 50001, portT: 50002},
	{hostname: "Electrum.hsmiths.com", portS: 8080, portT: 995},
	{hostname: "electrum3.hachre.de", portS: 50001, portT: 50002},
	{hostname: "b.1209k.com", portS: 50001, portT: 50002},
	{hostname: "elec.luggs.co", portS: 443, portT: 50002},
	{hostname: "btc.smsys.me", portS: 110, portT: 995},
}
func GetDefaultAvailableServer() ServerAddress{
	var conn net.Conn = nil
	var err error = nil
	rand.Seed(time.Now().Unix())
	var defaultServer ServerAddress
	for conn == nil {
		defaultServer= DefaultServers[rand.Intn(len(DefaultServers))]
		conn, err = net.DialTimeout("tcp",fmt.Sprintf("%v:%v",defaultServer.hostname,defaultServer.portS) , ConnectionTimeout)
	}
	if err != nil {
		fmt.Println("CANNOT ACCESS ANY OF DEFAULT SERVERS !!! CHECK CONNECTION")
		os.Exit(0)
	}
	return defaultServer
}

// FindElectrumServersIRC finds nodes to connect to by connecting to the Freenode #electrum channel.
func FindElectrumPublicServers() ([]string, error) {
	defaultServer:= GetDefaultAvailableServer()
	peerChan, err := GetAvailablePeers(defaultServer)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for hash := range peerChan {
			log.Printf("Address peer hash: %+v", hash)
		}
	}()

	
	if err != nil {
		log.Fatal("JRPC error:", err)
	}
	fmt.Println(fmt.Sprintf("Electrum default server: %v",defaultServer.hostname ))
	return []string{""},nil
}

func GetAvailablePeers(server ServerAddress) (<-chan string, error)  {
	node := NewNode()
	if err := node.ConnectTCP(fmt.Sprintf("%v:%v",server.hostname,server.portS)); err != nil {
		log.Fatal(err)
	}

	resp := &basicResp{}
	err := node.request("server.peers.subscribe", []string{}, resp)
	if err != nil {
		return nil, err
	}
	addressChan := make(chan string, 1)
	if len(resp.Result) > 0 {
		addressChan <- resp.Result
	}
	go func() {
		for msg := range node.listenPush("server.peers.subscribe") {
			resp := &struct {
				Params []string `json:"params"`
			}{}
			if err := json.Unmarshal(msg, resp); err != nil {
				log.Printf("ERR %s", err)
				return
			}
			if len(resp.Params) != 2 {
				log.Printf("address subscription params len != 2 %+v", resp.Params)
				continue
			}

				addressChan <- resp.Params[1]
		}
	}()
	return addressChan, err
}
