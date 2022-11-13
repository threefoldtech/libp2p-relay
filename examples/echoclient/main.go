package main

import (
	"context"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/threefoldtech/libp2p-relay/client"
)

const Protocol = "/echo/1.0.0"

func sayHello(h host.Host, peerID peer.ID) (err error) {
	log.Println("Trying to say hello to", peerID.Pretty())
	s, err := h.NewStream(context.Background(), peerID, Protocol)
	if err != nil {
		return
	}
	log.Println("sender saying hello")
	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		return
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		return
	}

	log.Printf("read reply: %q\n", out)
	s.Close()
	return
}
func main() {

	var hexPSK string
	var relay string
	var remotePeerID string
	flag.StringVar(&hexPSK, "psk", "", "32 bytes network PSK in hex")
	flag.StringVar(&relay, "relay", "", "relay multi-address")
	flag.StringVar(&remotePeerID, "remote", "", "Peer ID to connect to")
	flag.Parse()
	if hexPSK == "" {
		flag.Usage()
		log.Fatalln("The psk flag is required")
	}
	if remotePeerID == "" {
		flag.Usage()
		log.Fatalln("The remote flag is required")
	}
	targetID, err := peer.Decode(remotePeerID)
	if err != nil {
		log.Fatalln("Unable to hex decode the remote", err)
	}
	psk, err := hex.DecodeString(hexPSK)
	if err != nil {
		log.Fatalln("Unable to hex decode the PSK", err)
	}
	if len(psk) != 32 {
		log.Fatalln("The PSK should be 32 bytes")
	}

	relayAddrInfo, err := peer.AddrInfoFromString(relay)
	if err != nil {
		log.Fatalln(err)
	}

	libp2pctx := context.Background()
	p2pHost, peerRouting, err := client.CreateLibp2pHost(libp2pctx, 0, false, psk, nil, []peer.AddrInfo{*relayAddrInfo})
	if err != nil {
		panic(err)
	}
	log.Println("Started libp2p host on", p2pHost.Addrs())

	err = p2pHost.Connect(libp2pctx, *relayAddrInfo)
	if err != nil {
		log.Fatalln(err)
	}
	for {

		log.Println("Peers:", p2pHost.Peerstore().Peers())
		if err = client.ConnectToPeer(libp2pctx, p2pHost, peerRouting, relayAddrInfo, targetID); err != nil {
			log.Println("Unable to connect to remote:", err)
		}
		err = sayHello(p2pHost, targetID)
		if err != nil {
			log.Println("ERROR saying hello", err)
		}
		time.Sleep(time.Second * 10)

	}
}
