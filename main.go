package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"github.com/libp2p/go-libp2p/core/crypto"
)

var Version = "development"

func main() {
	var hexPSK string
	var hexPrivateKey string
	var tcpPort int
	var wsPort int
	var verbose bool
	var version bool

	flag.StringVar(&hexPSK, "psk", "", "32 bytes network preshared Key in hex")
	flag.StringVar(&hexPrivateKey, "idkey", "", "32 byte private key of the p2p Identity, if not provided, a random ID is generated")
	flag.IntVar(&tcpPort, "port", 0, "TCP port to listen on, if not set, a random port is taken")
	flag.IntVar(&wsPort, "wsport", -1, "websocket port to listen on, if not set, websockets are disabled")
	flag.BoolVar(&verbose, "verbose", false, "enable libp2p debug logging")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.Parse()
	if version {
		fmt.Println(Version)
		return
	}
	if hexPSK == "" {
		flag.Usage()
		log.Fatalln("The psk flag is required")
	}
	psk, err := hex.DecodeString(hexPSK)
	if err != nil {
		log.Fatalln("Unable to hex decode the psk", err)
	}
	if len(psk) != 32 {
		log.Fatalln("The PSK should be 32 bytes")
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

	libp2pctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	host, _, err := CreateLibp2pHost(libp2pctx, tcpPort, wsPort, psk, privKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("I am %s\n", host.ID())
	fmt.Printf("Public Addresses:\n")
	for _, addr := range host.Addrs() {
		fmt.Printf("\t%s/p2p/%s\n", addr, host.ID())
	}
	//Set up the host as a relay
	r, err := SetupRelay(host)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	for {
		fmt.Println("Peers:", host.Peerstore().Peers())
		time.Sleep(time.Minute * 2)
	}
}
