// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package goph

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	log.Printf("attempting connection to: %s", c.HostAddr)

	if c.Proto == "" {
		if strings.Contains(c.HostAddr, ".i2p") {
			log.Println("I2P address detected")
			c.Proto = darkssh.STREAMING
		} else if strings.Contains(c.HostAddr, ".onion") {
			log.Println("Tor address detected")
			c.Proto = darkssh.STREAMING
		} else {
			c.Proto = TCP
		}
	}

	if strings.Contains(c.HostAddr, ".i2p") {
		c.Conn, err = darkssh.DialI2P(c.Proto, fmt.Sprintf("%s", c.HostAddr), cfg)
	} else if strings.Contains(c.HostAddr, ".onion") {
		c.Conn, err = darkssh.DialTor(c.Proto, fmt.Sprintf("%s", c.HostAddr), cfg)
	} else {
		c.Conn, err = ssh.Dial(c.Proto, fmt.Sprintf("%s:%d", c.HostAddr, c.Port), cfg)
	}
	return
}

func (client *Client) Interact() error {

	fmt.Println("Welcome To darkssh B)")
	fmt.Printf("Connected to %s\n", client.HostAddr)
	fmt.Println("Type your shell command and enter.")
	fmt.Println("To download file from remote type: download remote/path local/path")
	fmt.Println("To upload file to remote type: upload local/path remote/path")
	fmt.Println("To exit type: exit")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")

	var (
		out   []byte
		err   error
		cmd   string
		parts []string
	)

loop:
	for scanner.Scan() {

		err = nil
		cmd = scanner.Text()
		parts = strings.Split(cmd, " ")

		if len(parts) < 1 {
			continue
		}

		switch parts[0] {

		case "exit":
			fmt.Println("goph bye!")
			break loop

		case "download":

			if len(parts) != 3 {
				fmt.Println("please type valid download command!")
				continue loop
			}

			err = client.Download(parts[1], parts[2])
			if err != nil {
				return err
			}

			fmt.Println("download err: ", err)
			break

		case "upload":

			if len(parts) != 3 {
				fmt.Println("please type valid upload command!")
				continue loop
			}

			err = client.Upload(parts[1], parts[2])

			fmt.Println("upload err: ", err)
			break

		default:

			out, err = client.Run(cmd)
			if err != nil {
				return err
			}
			fmt.Println(string(out), err)
		}

		fmt.Print("> ")
	}
	return nil
}
