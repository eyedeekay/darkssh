package main

import (
	"log"

	"github.com/eyedeekay/darkssh"
	"github.com/eyedeekay/goSam"
)

func main() {

	server, err := darkssh.Server(nil)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := darkssh.ListenI2P("st", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("starting ssh server on:", listener.(*goSam.Client).Base32()+".b32.i2p")
	log.Fatal(server.Serve(listener))
}
