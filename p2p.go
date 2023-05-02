package main

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"

	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

// CreateLibp2pHost creates a libp2p host with a dht in server mode to the bootstrap nodes
// If privateKey is nil, a libp2p host is started without a predefined peerID
// the tcpPort and the wsPort (websockets) can not be the same
func CreateLibp2pHost(ctx context.Context, tcpPort, wsport int, psk []byte, libp2pPrivKey crypto.PrivKey) (p2phost host.Host, peerRouting routing.PeerRouting, err error) {

	var idht *dht.IpfsDHT
	options := make([]libp2p.Option, 0, 0)

	// listen addresses

	// regular tcp connections
	options = append(options, libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", tcpPort),
	))

	// websocket connections
	if wsport >= 0 {
		options = append(options, libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/ws", wsport),
		))
	}

	// support TLS connections
	//options = append(options,
	//	libp2p.Security(libp2ptls.ID, libp2ptls.New))

	//Configure private network
	options = append(options, libp2p.PrivateNetwork(psk))

	if libp2pPrivKey != nil {
		options = append(options, libp2p.Identity(libp2pPrivKey))
	}

	//Explicitely set the transports to disable quic since it does not support private networks
	options = append(options, libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	))

	// Let's prevent our peer from having too many
	// connections by attaching a connection manager.
	cmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
	)
	if err != nil {
		return
	}
	options = append(options,
		libp2p.ConnectionManager(cmgr))

	// Attempt to open ports using uPNP for NATed hosts.
	// options = append(options,
	// 	libp2p.NATPortMap())

	// Enable the DHT as server
	options = append(options,
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h, dht.Mode(dht.ModeServer))
			return idht, err
		}))

	libp2phost, err := libp2p.New(options...)
	log.Println("Libp2p host started with PeerID", libp2phost.ID())

	return libp2phost, idht, err
}

func SetupRelay(p2phost host.Host) (r *relay.Relay, err error) {
	r, err = relay.New(p2phost)

	return
}
