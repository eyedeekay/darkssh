// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package goph

import (
	"fmt"
	"log"
	"strings"

	"github.com/eyedeekay/darkssh"
	"golang.org/x/crypto/ssh"
)

const (
	UDP string = "udp"
	TCP string = "tcp"
)

// Set new net connection to a client.
func Conn(c *Client, cfg *ssh.ClientConfig) (err error) {

	if c.Port == 0 {
		c.Port = 22
	}
	log.Printf("attempting connection to: %s", c.Addr)

	if c.Proto == "" {
		if strings.Contains(c.Addr, ".i2p") {
			log.Println("I2P address detected")
			c.Proto = darkssh.STREAMING
		} else if strings.Contains(c.Addr, ".onion") {
			log.Println("I2P address detected")
			c.Proto = darkssh.STREAMING
		} else {
			c.Proto = TCP
		}
	}

	if strings.Contains(c.Addr, ".i2p") {
		c.Conn, err = darkssh.DialI2P(c.Proto, fmt.Sprintf("%s", c.Addr), cfg)
	} else if strings.Contains(c.Addr, ".onion") {
		c.Conn, err = darkssh.DialTor(c.Proto, fmt.Sprintf("%s", c.Addr), cfg)
	} else {
		c.Conn, err = ssh.Dial(c.Proto, fmt.Sprintf("%s:%d", c.Addr, c.Port), cfg)
	}
	return
}
