package darkssh

import (
	"fmt"
	"github.com/eyedeekay/sam3"
	. "golang.org/x/crypto/ssh"
	"net"
)

var SAMAddress = "127.0.0.1:7656"

const (
  STREAMING string = "st"
  DATAGRAMS string = "dg"
)
// Dial returns an ssh.Client configured to connect via I2P. It accepts
// "st" or "dg" in the "Network" parameter, for "streaming" or "datagram"
// based connections. It is otherwise identical to ssh.Dial
func Dial(network, addr string, config *ClientConfig) (*Client, error) {
	switch network {
	case "st":
		conn, err := dialI2PStreaming(addr)
		if err != nil {
			return nil, err
		}
		c, chans, reqs, err := NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return NewClient(c, chans, reqs), nil
	case "dg":
		conn, err := dialI2PDatagrams(addr)
		if err != nil {
			return nil, err
		}
		c, chans, reqs, err := NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return NewClient(c, chans, reqs), nil
	default:
		return Dial("st", addr, config)
	}
}

func dialI2PStreaming(addr string) (net.Conn, error) {
	sam, err := sam3.NewSAM(SAMAddress)
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	keys, err := sam.NewKeys()
	if err != nil {
		return nil, err
	}
	fmt.Println("\tBuilding tunnel")
	ss, err := sam.NewStreamSession("streamTun", keys, []string{"inbound.length=2", "outbound.length=2", "inbound.lengthVariance=0", "outbound.lengthVariance=0", "inbound.quantity=3", "outbound.quantity=3"})
	if err != nil {
		return nil, err
	}
	forumAddr, err := sam.Lookup(addr)
	if err != nil {
		return nil, err
	}
	return ss.DialI2P(forumAddr)
}

func dialI2PDatagrams(addr string) (net.Conn, error) {
	sam, err := sam3.NewSAM(SAMAddress)
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	keys, err := sam.NewKeys()
	if err != nil {
		return nil, err
	}
	fmt.Println("\tBuilding tunnel")
	return sam.NewDatagramSession("streamTun", keys, []string{"inbound.length=2", "outbound.length=2", "inbound.lengthVariance=0", "outbound.lengthVariance=0", "inbound.quantity=3", "outbound.quantity=3"}, 0)
}
