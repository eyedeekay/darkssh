package darkssh

import (
	"log"
	"net"

	"github.com/cretz/bine/tor"
	"golang.org/x/crypto/ssh"
)

var TORHost = "127.0.0.1"
var SOCKSPort = "9050"
var CONTROLPort = "9051"

func SOCKSAddress() string {
	return TORHost + ":" + SOCKSPort
}

func CONTROLAddress() string {
	return TORHost + ":" + SOCKSPort
}

const (
	TORTCP string = "tor"
)

// Dial returns an ssh.Client configured to connect via I2P. It accepts
// "st" or "dg" in the "Network" parameter, for "streaming" or "datagram"
// based connections. It is otherwise identical to ssh.Dial
func DialTor(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	switch network {
	case "tor":
		conn, err := dialTorStreaming(addr)
		if err != nil {
			return nil, err
		}
		c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return ssh.NewClient(c, chans, reqs), nil
	default:
		return Dial("tor", addr, config)
	}
}

func dialTorStreaming(addr string) (net.Conn, error) {
	log.Println("\tBuilding connection")
	t, err := tor.Start(nil, nil)
	if err != nil {
		return nil, err
	}
	d, err := t.Dialer(nil, &tor.DialConf{ProxyAddress: SOCKSAddress()})
	if err != nil {
		return nil, err
	}
	return d.DialContext(nil, "tcp", addr)
}
