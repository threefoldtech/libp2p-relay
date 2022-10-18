package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"
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
	libp2pctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p2pHost, _, err := CreateLibp2pHost(libp2pctx, tcpPort, psk, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Started libp2p host on", p2pHost.Addrs())
	//Set up the host as a relay
	r, err := SetupRelay(p2pHost)
	defer r.Close()
	for {
		fmt.Println("My Peers:", p2pHost.Peerstore().Peers())
		time.Sleep(time.Minute * 2)
	}
}
