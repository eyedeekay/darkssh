package darkssh

import (
	"io"
	"log"
	"net"

	"github.com/eyedeekay/goSam"
	"github.com/eyedeekay/sam3"
	sshd "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
)

// SAMHost is the SAM API bridge host
var SAMHost = "127.0.0.1"

// SAMPort is the SAM API bridge port
var SAMPort = "7656"

// SAMHostAddress combines SAMHost and SAMPort
func SAMHostAddress() string {
	return SAMHost + ":" + SAMPort
}

const (
	// STREAMING is an I2P Streaming Session
	STREAMING string = "st"
	// DATAGRAMS is an I2P Datagram Session
	DATAGRAMS string = "dg"
)

// DialI2P returns an ssh.Client configured to connect via I2P. It accepts
// "st" or "dg" in the "Network" parameter, for "streaming" or "datagram"
// based connections. It is otherwise identical to ssh.Dial
func DialI2P(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	switch network {
	case "st":
		conn, err := DialI2PStreaming(addr)
		if err != nil {
			return nil, err
		}
		c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return ssh.NewClient(c, chans, reqs), nil
	case "dg":
		conn, err := DialI2PDatagrams(addr)
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

func DialI2PStreaming(addr string) (net.Conn, error) {
	sam, err := goSam.NewClientFromOptions(
		goSam.SetHost(SAMHost),
		goSam.SetPort(SAMPort),
		goSam.SetDebug(true),
		goSam.SetInLength(2),
		goSam.SetOutLength(2),
		goSam.SetInQuantity(3),
		goSam.SetOutQuantity(3),
		goSam.SetInBackups(2),
		goSam.SetOutBackups(2),
		goSam.SetCloseIdle(false),
		goSam.SetReduceIdle(false),
	)
	if err != nil {
		return nil, err
	}
	log.Println("\tBuilding tunnel")
	return sam.Dial("tcp", addr)
}

func DialI2PDatagrams(addr string) (net.Conn, error) {
	sam, err := sam3.NewSAM(SAMHostAddress())
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	keys, err := sam.NewKeys()
	if err != nil {
		return nil, err
	}
	log.Println("\tBuilding tunnel")
	return sam.NewDatagramSession("streamTun", keys, Options_SSH, 0)
}

var Options_SSH = []string{"inbound.length=2", "outbound.length=2", "inbound.lengthVariance=0", "outbound.lengthVariance=0", "inbound.quantity=3", "outbound.quantity=3", "inbound.backupQuantity=2", "outbound.backupQuantity=2", "i2cp.closeOnIdle=false", "i2cp.reduceOnIdle=false", "i2cp.leaseSetEncType=4,0"}

func Server(config *sshd.Option) (*sshd.Server, error) {
	forwardHandler := &sshd.ForwardedTCPHandler{}
	server := sshd.Server{
		LocalPortForwardingCallback: sshd.LocalPortForwardingCallback(func(ctx sshd.Context, dhost string, dport uint32) bool {
			log.Println("Accepted forward", dhost, dport)
			return true
		}),
		Handler: sshd.Handler(func(s sshd.Session) {
			io.WriteString(s, "Remote forwarding available...\n")
			select {}
		}),
		ReversePortForwardingCallback: sshd.ReversePortForwardingCallback(func(ctx sshd.Context, host string, port uint32) bool {
			log.Println("attempt to bind", host, port, "granted")
			return true
		}),
		RequestHandlers: map[string]sshd.RequestHandler{
			"tcpip-forward":        forwardHandler.HandleSSHRequest,
			"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
		},
	}
	//	server.SetOption(config)
	return &server, nil
}

func ListenI2P(network string, config *sshd.Option) (net.Listener, error) {
	switch network {
	case "st":
		return ListenI2PStreaming()
	case "dg":
		return ListenI2PDatagrams()
	default:
		return ListenI2P("st", config)
	}
}

func ListenI2PStreaming() (net.Listener, error) {
	sam, err := goSam.NewClientFromOptions(
		goSam.SetHost(SAMHost),
		goSam.SetPort(SAMPort),
		goSam.SetUnpublished(false),
		goSam.SetDebug(true),
		goSam.SetInLength(2),
		goSam.SetOutLength(2),
		goSam.SetInQuantity(3),
		goSam.SetOutQuantity(3),
		goSam.SetInBackups(2),
		goSam.SetOutBackups(2),
		goSam.SetCloseIdle(false),
		goSam.SetReduceIdle(false),
	)
	if err != nil {
		return nil, err
	}
	log.Println("\tBuilding tunnel")
	return sam.Listen()
}

func ListenI2PDatagrams() (net.Listener, error) {
	sam, err := sam3.NewSAM(SAMHostAddress())
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	keys, err := sam.NewKeys()
	if err != nil {
		return nil, err
	}
	log.Println("\tBuilding tunnel")
	return sam.NewDatagramSession("DGTUN", keys, Options_SSH, 0)
}
