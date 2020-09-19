package darkssh

import (
	"log"
	"net"

	"github.com/cretz/bine/tor"
	"golang.org/x/crypto/ssh"
)

// TORHost is the host where Tor is running
var TORHost = "127.0.0.1"

// SOCKSPort is the port used for the Tor SOCKS proxy
var SOCKSPort = "9050"

// CONTROLPort is the port used for controlling Tor
var CONTROLPort = "9051"

// SOCKSHostAddress gives you the address of the Tor SOCKS port
func SOCKSHostAddress() string {
	return TORHost + ":" + SOCKSPort
}

// CONTROLHostAddress gets you the address of the Tor Control Port
func CONTROLHostAddress() string {
	return TORHost + ":" + SOCKSPort
}

const (
	// TORTCP a TOR TCP session
	TORTCP string = "tor"
)

// DialTor returns an ssh.Client configured to connect via Tor. It accepts
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
		return DialTor("tor", addr, config)
	}
}

func dialTorStreaming(addr string) (net.Conn, error) {
	log.Println("\tBuilding connection")
	t, err := tor.Start(nil, nil)
	if err != nil {
		return nil, err
	}
	d, err := t.Dialer(nil, &tor.DialConf{ProxyAddress: SOCKSHostAddress()})
	if err != nil {
		return nil, err
	}
	return d.DialContext(nil, "tcp", addr)
}
