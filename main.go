package main

import (
	"context"
	"encoding/hex"
	"flag"
	"log"
	"time"

	"github.com/threefoldtech/libp2p-relay/communication"
)

func main() {
	var hexPSK string
	var tcpPort int
	flag.StringVar(&hexPSK, "psk", "", "32 bytes network PSK in hex")
	flag.IntVar(&tcpPort, "port", 0, "TCP port to listen on, if not set, a random port is taken")
	flag.Parse()
	if hexPSK == "" {
		flag.Usage()
		log.Fatalln("The psk flag is required")
	}
	psk, err := hex.DecodeString(hexPSK)
	if err != nil {
		log.Fatalln("Unable to hex decode the PSK", err)
	}
	if len(psk) != 32 {
		log.Fatalln("The PSK should be 32 bytes")
	}
	libp2pctx := context.Background()
	host, _, err := communication.CreateLibp2pHost(libp2pctx, tcpPort, psk)
	if err != nil {
		panic(err)
	}
	log.Println("Started libp2p host on", host.Addrs())

	for {
		time.Sleep(time.Second * 10)
	}
}
