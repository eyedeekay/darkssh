// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package goph

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	//"sync/atomic"
	"time"

	"github.com/eyedeekay/darkssh"
	"golang.org/x/crypto/ssh"
)

func (t Client) BindTunnel(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var once sync.Once // Only print errors once per session
		func() {
			// Connect to the server host via SSH.
			/*
				&ssh.ClientConfig{
					User:            t.User,
					Auth:            t.Auth,
					HostKeyCallback: t.HostKeys,
					Timeout:         300 * time.Second,
				}
			*/
			cl, err := NewConn(t.User, t.HostAddr, t.Auth, t.HostKeys)
			if err != nil {
				once.Do(func() { fmt.Printf("(%v) SSH dial error: %v", t, err) })
				return
			}
			wg.Add(1)
			go t.KeepAliveMonitor(&once, wg, cl.Interactive.Conn)
			defer cl.Close()

			// Attempt to bind to the inbound socket.
			var ln net.Listener
			switch t.Mode {
			case '>':
				ln, err = net.Listen("tcp", t.BindAddr)
			case '<':
				ln, err = cl.Interactive.Conn.Listen("tcp", t.BindAddr)
			}
			if err != nil {
				once.Do(func() { fmt.Printf("(%v) bind error: %v", t, err) })
				return
			}

			// The socket is binded. Make sure we close it eventually.
			bindCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			go func() {
				cl.Interactive.Conn.Wait()
				cancel()
			}()
			go func() {
				<-bindCtx.Done()
				once.Do(func() {}) // Suppress future errors
				ln.Close()
			}()

			fmt.Printf("(%v) binded tunnel", t)
			defer fmt.Printf("(%v) collapsed tunnel", t)

			// Accept all incoming connections.
			for {
				cn1, err := ln.Accept()
				if err != nil {
					once.Do(func() { fmt.Printf("(%v) accept error: %v", t, err) })
					return
				}
				wg.Add(1)
				go t.DialTunnel(bindCtx, wg, cl.Interactive.Conn, cn1)
			}
		}()

		select {
		case <-ctx.Done():
			return
		case <-time.After(t.RetryInterval):
			fmt.Printf("(%v) retrying...", t)
		}
	}
}

func (t Client) DialTunnel(ctx context.Context, wg *sync.WaitGroup, client *ssh.Client, cn1 net.Conn) {
	defer wg.Done()

	// The inbound connection is established. Make sure we close it eventually.
	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-connCtx.Done()
		cn1.Close()
	}()

	// Establish the outbound connection.
	var cn2 net.Conn
	var err error
	switch t.Mode {
	case '>':
		cn2, err = client.Dial("tcp", t.DialAddr)
	case '<':
		cn2, err = darkssh.DialI2PStreaming("tcp", t.DialAddr)
	}
	if err != nil {
		fmt.Printf("(%v) dial error: %v", t, err)
		return
	}

	go func() {
		<-connCtx.Done()
		cn2.Close()
	}()

	fmt.Printf("(%v) connection established", t)
	defer fmt.Printf("(%v) connection closed", t)

	// Copy bytes from one connection to the other until one side closes.
	var once sync.Once
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go func() {
		defer wg2.Done()
		defer cancel()
		if _, err := io.Copy(cn1, cn2); err != nil {
			once.Do(func() { fmt.Printf("(%v) connection error: %v", t, err) })
		}
		once.Do(func() {}) // Suppress future errors
	}()
	go func() {
		defer wg2.Done()
		defer cancel()
		if _, err := io.Copy(cn2, cn1); err != nil {
			once.Do(func() { fmt.Printf("(%v) connection error: %v", t, err) })
		}
		once.Do(func() {}) // Suppress future errors
	}()
	wg2.Wait()
}
