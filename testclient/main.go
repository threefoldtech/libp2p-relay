package main

import (
	"context"
	"encoding/hex"
	"flag"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/threefoldtech/libp2p-relay/communication"
)

func main() {

	var hexPSK string
	var relay string
	var remotePeerID string
	flag.StringVar(&hexPSK, "psk", "", "32 bytes network PSK in hex")
	flag.StringVar(&relay, "relay", "", "relay libp2p address")
	flag.StringVar(&remotePeerID, "remote", "", "Peer ID to connect to")
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
	host, peerRouting, err := communication.CreateLibp2pHost(libp2pctx, "", psk)
	if err != nil {
		panic(err)
	}
	log.Println("Started libp2p host on", host.Addrs())
	pi, err := peer.AddrInfoFromString(relay)
	if err != nil {
		log.Fatalln(err)
	}

	err = host.Connect(libp2pctx, *pi)
	if err != nil {
		log.Fatalln(err)
	}
	for {

		if remotePeerID != "" {
			if err = communication.ConnectToPeer(libp2pctx, host, peerRouting, peer.ID(remotePeerID)); err != nil {
				log.Println("Unable to connect to remote", err)
			}
		}
		log.Println("Peers:", host.Peerstore().Peers())
		time.Sleep(time.Second * 10)

	}
}
