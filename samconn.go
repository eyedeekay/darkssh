package darkssh

import (
	"log"
	"net"

	"github.com/eyedeekay/gosam"
	"github.com/eyedeekay/sam3"
	"golang.org/x/crypto/ssh"
)

var SAMHost = "127.0.0.1"
var SAMPort = "7656"

func SAMAddress() string {
	return SAMHost + ":" + SAMPort
}

const (
	STREAMING string = "st"
	DATAGRAMS string = "dg"
)

// Dial returns an ssh.Client configured to connect via I2P. It accepts
// "st" or "dg" in the "Network" parameter, for "streaming" or "datagram"
// based connections. It is otherwise identical to ssh.Dial
func DialI2P(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	switch network {
	case "st":
		conn, err := dialI2PStreaming(addr)
		if err != nil {
			return nil, err
		}
		c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return ssh.NewClient(c, chans, reqs), nil
	case "dg":
		conn, err := dialI2PDatagrams(addr)
		if err != nil {
			return nil, err
		}
		c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return ssh.NewClient(c, chans, reqs), nil
	default:
		return DialI2P("st", addr, config)
	}
}

func dialI2PStreaming(addr string) (net.Conn, error) {
	sam, err := goSam.NewClientFromOptions(
		goSam.SetHost(SAMHost),
		goSam.SetPort(SAMPort),
		goSam.SetDebug(false),
	)
	if err != nil {
		return nil, err
	}
	log.Println("\tBuilding tunnel")
	return sam.Dial("tcp", addr)
}

func dialI2PDatagrams(addr string) (net.Conn, error) {
	sam, err := sam3.NewSAM(SAMAddress())
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	keys, err := sam.NewKeys()
	if err != nil {
		return nil, err
	}
	log.Println("\tBuilding tunnel")
	return sam.NewDatagramSession("streamTun", keys, []string{"inbound.length=2", "outbound.length=2", "inbound.lengthVariance=0", "outbound.lengthVariance=0", "inbound.quantity=3", "outbound.quantity=3"}, 0)
}
