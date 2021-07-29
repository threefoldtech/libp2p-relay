package communication

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p-core/host"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
)

//CreateLibp2pHost creates a libp2p host with peerRouting enabled and connects to the bootstrap nodes
// If privateKey is nil, a libp2p host is started without a predefined peerID
func CreateLibp2pHost(ctx context.Context, tcpPort, quicUdpPort string, psk []byte) (host.Host, *dht.IpfsDHT, error) {

	if tcpPort == "" {
		tcpPort = "0"
	}
	if quicUdpPort == "" {
		quicUdpPort = "0"
	}
	var idht *dht.IpfsDHT
	var err error
	options := make([]libp2p.Option, 0, 0)
	// Multiple listen addresses
	options = append(options, libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/"+tcpPort,             // regular tcp connections
		"/ip4/0.0.0.0/udp/"+quicUdpPort+"/quic", // a UDP endpoint for the QUIC transport
	))
	// support TLS connections
	options = append(options,
		libp2p.Security(libp2ptls.ID, libp2ptls.New))

	//Configure private network
	options = append(options, libp2p.PrivateNetwork(psk))

	// support any other default transports (TCP)
	options = append(options,
		libp2p.DefaultTransports)

	// Let's prevent our peer from having too many
	// connections by attaching a connection manager.
	options = append(options,
		libp2p.ConnectionManager(connmgr.NewConnManager(
			100,         // Lowwater
			400,         // HighWater,
			time.Minute, // GracePeriod
		)))
	// Attempt to open ports using uPNP for NATed hosts.
	options = append(options,
		libp2p.NATPortMap())
	// Let this host use the DHT to find other hosts
	options = append(options,
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}))
	// Let this host use relays and advertise itself on relays if
	// it finds it is behind NAT. Use libp2p.Relay(options...) to
	// enable active relays and more.
	options = append(options, libp2p.EnableAutoRelay())

	libp2phost, err := libp2p.New(ctx, options...)
	log.Println("Libp2p host started with PeerID", libp2phost.ID())

	return libp2phost, idht, err
}
func ConnectToPeer(ctx context.Context, p2phost host.Host, hostRouting routing.PeerRouting, peerID peer.ID) (err error) {

	findPeerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	peeraddrInfo, err := hostRouting.FindPeer(findPeerCtx, peerID)
	if err != nil {
		return
	}
	ConnectPeerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	err = p2phost.Connect(ConnectPeerCtx, peeraddrInfo)
	return
}
