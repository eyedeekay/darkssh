package darkssh

import (
	//"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strings"
)

func DialConn(network, addr string) (net.Conn, error) {
	if strings.Contains(addr, ".i2p") {
		log.Println("I2P address detected")
		return DialI2PStreaming(network, addr)
	} else if strings.Contains(addr, ".onion") {
		log.Println("Tor address detected")
		return DialTorStreaming(network, addr)
	}
	return net.Dial(network, addr)
}
