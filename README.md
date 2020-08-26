DarkSSH - SSH-over-Anonymous-Networks for Go
============================================

This is a tool for automatically connecting to SSH Servers and SSH-based
Services which are hosted on `.i2p` and `.onion` resources. Given an .i2p
or .onion host, it will dial out using their respective sockets. It is
used to make managing many `.i2p` or `.onion`-hosted SSH services easier
by natively handling the `known_hosts` file and automatically handling
proxy setup and client key-management for services which use blocklisting
and allowlisting facilities provided by the hidden services, or even more
sophisticated features like Encrypted LeaseSets. As a fringe benefit,
when addressing services by their cryptographic identifiers(i.e. the 
`.b32.i2p` or `.onion` domains) there is no chance of impersonation.
Eventually, it will implement a drop-in replacement for a real SSH client
so it can be used as a ProxyCommand or as part of a `.i2p` or `.onion` only
selfhosted workflow.

What's in this repository:

```bash        
# a terminal SSH client - interface is *UNSTABLE*, forked from goph
# for modification
./cmd/darkssh
# a slightly-modified version of melbahja/goph, which automatically
# configures itself for I2P and Tor Transports
./goph
# implementations of the required interfaces for x/crypto/ssh
./
```

The goal is to be exactly compatible with any other SSH client, so things
that proxy commands to SSH, like rsync or SSHFS, can use it instead when
someone wants to use such a tool over Tor or I2P.

Eventually, an SSH server will also be implemented.
