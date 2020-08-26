// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package goph

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"strings"
)

type Auth []ssh.AuthMethod

// Get auth method from raw password.
func Password(pass string) Auth {
	return Auth{
		ssh.Password(pass),
	}
}

func askPass(msg string) string {

	fmt.Print(msg)

	pass, err := terminal.ReadPassword(0)

	if err != nil {
		panic(err)
	}

	fmt.Println("")

	return strings.TrimSpace(string(pass))
}

func getPassphrase(ask bool) string {

	if ask {
		return askPass("Enter Private Key Passphrase: ")
	}

	return ""
}

// Get auth method from private key with or without passphrase.
func Key(prvFile string, passphrase string) Auth {

	var err error
	signer, err := GetSigner(prvFile, passphrase)

	if err != nil {
		passphrase := getPassphrase(true)
		signer, err = GetSigner(prvFile, passphrase)
		if err != nil {
			panic(err)
		}
	}

	return Auth{
		ssh.PublicKeys(signer),
	}
}

// Get private key signer.
func GetSigner(prvFile string, passphrase string) (ssh.Signer, error) {

	var (
		err    error
		signer ssh.Signer
	)

	privateKey, err := ioutil.ReadFile(prvFile)

	if err != nil {

		return nil, err

	} else if passphrase != "" {

		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))

	} else {

		signer, err = ssh.ParsePrivateKey(privateKey)
	}

	return signer, err
}
