// Package electrum is a electrum bitcoin client.
package electrum

import (
	"fmt"
	"math/rand"
	"time"
	"net"
	"os"
)

const (
	CONNECTION_TIMEOUT = 5 * time.Second
)
type ServerAddress struct {
	hostname string
	portT int
	portS int
}
var DEFAULT_SERVERS = []ServerAddress{
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
	defaultServer:= DEFAULT_SERVERS[rand.Intn(len(DEFAULT_SERVERS))]
	for conn == nil {
		conn, err = net.DialTimeout("tcp",fmt.Sprintf("%v:%v",defaultServer.hostname,defaultServer.portT) , CONNECTION_TIMEOUT)
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
	fmt.Println(fmt.Sprintf("Electrum default server: %v",defaultServer.hostname ))
	return []string{""},nil
}

