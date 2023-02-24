package client

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

// CreateLibp2pHost creates a libp2p host with a dht in server mode to the bootstrap nodes
// listen idicates wether or not a tcpport should be opened for the host to listen on.
// If privateKey is nil, a libp2p host is started without a predefined peerID
func CreateLibp2pHost(ctx context.Context, tcpPort int, listen bool, psk []byte, libp2pPrivKey crypto.PrivKey, relays []peer.AddrInfo) (p2phost host.Host, peerRouting routing.PeerRouting, err error) {
	var idht *dht.IpfsDHT
	options := make([]libp2p.Option, 0)
	// listen addresses
	if listen {
		options = append(options, libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", tcpPort), // regular tcp connections
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
	options = append(options,
		libp2p.NATPortMap())
	// Enable the DHT
	options = append(options,
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h, dht.BootstrapPeers(relays...), dht.Mode(dht.ModeAuto))
			return idht, err
		}))
	// Let this host use relays and advertise itself on relays if
	// it finds it is behind NAT.
	options = append(options, libp2p.EnableAutoRelayWithStaticRelays(relays))

	libp2phost, err := libp2p.New(options...)
	log.Println("Libp2p host started with PeerID", libp2phost.ID())

	return libp2phost, idht, err
}

func ConnectToPeer(ctx context.Context, p2phost host.Host, hostRouting routing.PeerRouting, relay *peer.AddrInfo, peerID peer.ID) (err error) {
	targetMA, e := multiaddr.NewMultiaddr("/p2p/" + relay.ID.String() + "/p2p-circuit/p2p/" + peerID.String())
	if e != nil {
		err = e
		return
	}
	peeraddrInfo := peer.AddrInfo{
		ID:    peerID,
		Addrs: []multiaddr.Multiaddr{targetMA},
	}

	ConnectPeerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	return p2phost.Connect(ConnectPeerCtx, peeraddrInfo)
}
