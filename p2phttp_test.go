package p2phttp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	gostream "github.com/hsanjuan/go-libp2p-gostream"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// newHost illustrates how to build a libp2p host with secio using
// a randomly generated key-pair
func newHost(t *testing.T, listen multiaddr.Multiaddr) host.Host {
	priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	ps := peerstore.NewPeerstore()
	err = ps.AddPubKey(pid, pub)
	if err != nil {
		t.Fatal(err)
	}
	err = ps.AddPrivKey(pid, priv)
	if err != nil {
		t.Fatal(err)
	}

	network, err := swarm.NewNetwork(
		context.Background(),
		[]multiaddr.Multiaddr{listen},
		pid,
		ps,
		nil)

	if err != nil {
		t.Fatal(err)
	}

	host := bhost.New(network)
	return host
}

func TestServerClient(t *testing.T) {
	m1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10000")
	m2, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10001")
	srvHost := newHost(t, m1)
	clientHost := newHost(t, m2)
	defer srvHost.Close()
	defer clientHost.Close()

	srvHost.Peerstore().AddAddrs(clientHost.ID(), clientHost.Addrs(), peerstore.PermanentAddrTTL)
	clientHost.Peerstore().AddAddrs(srvHost.ID(), srvHost.Addrs(), peerstore.PermanentAddrTTL)

	listener, err := gostream.Listen(srvHost, P2PProtocol)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hi!"))
		})
		server := &http.Server{}
		server.Serve(listener)
	}()

	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", NewTransport(clientHost))
	client := &http.Client{Transport: tr}
	res, err := client.Get(fmt.Sprintf("libp2p://%s/hello", srvHost.ID().Pretty()))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(text) != "Hi!" {
		t.Errorf("expected Hi! but got %s", text)
	}
}
