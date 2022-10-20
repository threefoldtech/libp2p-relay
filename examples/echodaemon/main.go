package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"flag"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/threefoldtech/libp2p-relay/client"
)

const Protocol = "/echo/1.0.0"

// doEcho reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}
func main() {

	var hexPSK string
	var relay string
	var listen bool
	flag.StringVar(&hexPSK, "psk", "", "32 bytes network PSK in hex")
	flag.StringVar(&relay, "relay", "", "relay multi-address")
	flag.BoolVar(&listen, "listen", true, "open a tcp port to listen on")
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

	relayAddrInfo, err := peer.AddrInfoFromString(relay)
	if err != nil {
		log.Fatalln(err)
	}
	libp2pctx := context.Background()
	p2pHost, _, err := client.CreateLibp2pHost(libp2pctx, 0, listen, psk, nil, []peer.AddrInfo{*relayAddrInfo})
	if err != nil {
		panic(err)
	}
	log.Println("Started libp2p host on", p2pHost.Addrs())

	// Set a stream handler on the host.
	p2pHost.SetStreamHandler(Protocol, func(s network.Stream) {
		log.Println("listener received new stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	err = p2pHost.Connect(libp2pctx, *relayAddrInfo)
	if err != nil {
		log.Fatalln(err)
	}
	for {

		log.Println("Peers:", p2pHost.Peerstore().Peers())
		time.Sleep(time.Minute)

	}
}
