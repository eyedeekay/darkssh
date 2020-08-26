DarkSSH - SSH-over-Anonymous-Networks for Go
============================================

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
