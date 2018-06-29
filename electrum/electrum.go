// Package electrum is a electrum bitcoin client.
package electrum

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"net"
	"os"
	"log"
	"encoding/json"
	"strconv"
)

const (
	ConnectionTimeout = 10 * time.Second
)

type ServerAddress struct {
	Hostname string
	PortT    int
	PortS    int
}

func (server *ServerAddress) GetAddressString() (string) {
	return fmt.Sprintf("%v:%v", server.Hostname, server.PortT)
}

var DefaultServers = []ServerAddress{
	{Hostname: "erbium1.sytes.net", PortT: 50001, PortS: 50002},
	{Hostname: "ecdsa.net", PortT: 50001, PortS: 110},
	{Hostname: "gh05.geekhosters.com", PortT: 50001, PortS: 50002},
	{Hostname: "eVPS.hsmiths.com", PortT: 50001, PortS: 50002},
	{Hostname: "electrum.anduck.net", PortT: 50001, PortS: 50002},
	{Hostname: "electrum.no-ip.org", PortT: 50001, PortS: 50002},
	{Hostname: "electrum.be", PortT: 50001, PortS: 50002},
	{Hostname: "helicarrier.bauerj.eu", PortT: 50001, PortS: 50002},
	{Hostname: "elex01.blackpole.online", PortT: 50001, PortS: 50002},
	{Hostname: "electrumx.not.fyi", PortT: 50001, PortS: 50002},
	{Hostname: "node.xbt.eu", PortT: 50001, PortS: 50002},
	{Hostname: "kirsche.emzy.de", PortT: 50001, PortS: 50002},
	{Hostname: "electrum.villocq.com", PortT: 50001, PortS: 50002},
	{Hostname: "us11.einfachmalnettsein.de", PortT: 50001, PortS: 50002},
	{Hostname: "electrum.trouth.net", PortT: 50001, PortS: 50002},
	{Hostname: "Electrum.hsmiths.com", PortT: 8080, PortS: 995},
	{Hostname: "electrum3.hachre.de", PortT: 50001, PortS: 50002},
	{Hostname: "b.1209k.com", PortT: 50001, PortS: 50002},
}

type ServerPeersResponse struct {
	Servers []ServerPeerResponse `json:"result"`
}

type ServerPeerResponse struct {
	IpAddress string
	Hostname  string
	Ports     []string
}

func (peer *ServerPeerResponse) ToServerAddress() (*ServerAddress) {
	var peerPortSSL = 0
	var peerPortTCP = 0
	if len(peer.Ports) > 1 {
		peerPortSSL, _ = strconv.Atoi(trimLeftChars(peer.Ports[1], 1))
		if len(peer.Ports) > 2 {
			peerPortTCP, _ = strconv.Atoi(trimLeftChars(peer.Ports[2], 1))
		}
	}
	return &ServerAddress{Hostname: peer.Hostname, PortS: peerPortSSL, PortT: peerPortTCP}
}
func trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

func (i *ServerPeerResponse) UnmarshalJSON(b []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	// Now grab the nested values using the right types for each
	var ipAddress string
	var hostname string
	var ports []string
	if err := json.Unmarshal(raw[0], &ipAddress); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[1], &hostname); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[2], &ports); err != nil {
		return err
	}

	*i = ServerPeerResponse{ipAddress, hostname, ports}

	return nil
}

func GetDefaultAvailableServer() ServerAddress {
	var conn net.Conn = nil
	var err error = nil
	rand.Seed(time.Now().Unix())
	var defaultServer ServerAddress
	for conn == nil {
		defaultServer = DefaultServers[rand.Intn(len(DefaultServers))]
		conn, err = net.DialTimeout("tcp", defaultServer.GetAddressString(), ConnectionTimeout)
	}
	if err != nil {
		fmt.Println("CANNOT ACCESS ANY OF DEFAULT SERVERS !!! CHECK CONNECTION")
		os.Exit(0)
	}
	defer conn.Close()
	return defaultServer
}

// FindElectrumServersIRC finds nodes to connect to by connecting to the Freenode #electrum channel.
func FindElectrumPublicServer() (*ServerAddress, error) {
	defaultServer := GetDefaultAvailableServer()
	servers, err := GetAvailablePeersCall(defaultServer)
	if err != nil {
		log.Fatal("JRPC error:", err)
	}
	var server = GetRandomPeerAddress(servers)
	return server, err
}
func FindElectrumPublicServers() (*[]ServerAddress, error) {
	defaultServer := GetDefaultAvailableServer()
	servers, err := GetAvailablePeersCall(defaultServer)
	if err != nil {
		log.Fatal("JRPC error:", err)
	}
	var serverAddresses = make([]ServerAddress, 0)
	for _, server := range servers {
		// Omit onion servers
		if !strings.HasSuffix(server.Hostname,".onion") {
			serverAddresses = append(serverAddresses, *server.ToServerAddress())
		}
	}
	return &serverAddresses, err
}
func GetRandomPeerAddress(peers []ServerPeerResponse) *ServerAddress {
	rand.Seed(time.Now().Unix())
	var peerServer ServerPeerResponse
	for len(peerServer.Ports) < 3 {
		peerServer = peers[rand.Intn(len(peers))]
	}

	return peerServer.ToServerAddress()
}
func GetAvailablePeersCall(server ServerAddress) ([]ServerPeerResponse, error) {
	node := NewNode()
	if err := node.ConnectTCP(server.GetAddressString()); err != nil {
		log.Fatal(err)
	}
	var result ServerPeersResponse
	err := node.request("server.peers.subscribe", []string{}, &result)
	return result.Servers, err
}

func GetAvailablePeersSubscribe(server ServerAddress) (<-chan string, error) {
	node := NewNode()
	if err := node.ConnectTCP(fmt.Sprintf("%v:%v", server.Hostname, server.PortS)); err != nil {
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
