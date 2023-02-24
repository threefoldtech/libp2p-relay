package main

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/threefoldtech/libp2p-relay/client"

	logging "github.com/ipfs/go-log/v2"
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
	var hexPrivateKey string
	var relay string
	var listen bool
	var verbose bool
	flag.StringVar(&hexPSK, "psk", "", "32 bytes network PSK in hex")
	flag.StringVar(&relay, "relay", "", "relay multi-address")
	flag.StringVar(&hexPrivateKey, "idkey", "", "32 byte private key of the p2p Identity, if not provided, a random ID is generated")
	flag.BoolVar(&listen, "listen", true, "open a tcp port to listen on")
	flag.BoolVar(&verbose, "verbose", false, "enable libp2p debug logging")
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

	var privKey crypto.PrivKey
	if hexPrivateKey != "" {
		privKeySeed, err := hex.DecodeString(hexPrivateKey)
		if err != nil {
			log.Fatalln("Unable to hex decode the idkey", err)
		}
		if len(privKeySeed) != 32 {
			log.Fatalln("The idKey should be 32 bytes")
		}
		privKey, err = crypto.UnmarshalEd25519PrivateKey(
			ed25519.NewKeyFromSeed(privKeySeed),
		)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if verbose {
		logging.SetDebugLogging()
	}

	libp2pctx := context.Background()
	p2pHost, _, err := client.CreateLibp2pHost(libp2pctx, 0, listen, psk, privKey, []peer.AddrInfo{*relayAddrInfo}, true)
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
