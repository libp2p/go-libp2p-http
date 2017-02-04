// Package p2phttp allows to serve HTTP endpoints and make HTTP requests through
// LibP2P (https://github.com/libp2p/libp2p) using Go's standard "http" and
// "net" stacks.
//
// Instead of the regular "host:port" addressing, `p2phttp` uses a Peer ID
// and lets LibP2P take care of the routing, thus taking advantage
// of features like multi-routes,  NAT transversal and stream multiplexing
// over a single connection.
//
// When already running a LibP2P facility, this package allows to expose
// existing HTTP-based services (like REST APIs) through LibP2P and to
// use those services with minimal changes to the code-base.
//
// For example, a simple http.Server on LibP2P works as:
//
//	listener, _ := p2phttp.Listen(host1)
//	defer listener.Close()
//	go func() {
//		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
//			w.Write([]byte("Hi!"))
//		})
//		server := &http.Server{}
//		server.Serve(listener)
//	}
//      ...
//
// As shown above, a Server only needs a custom listener which uses a libP2P
// host.
//
// On the other side, a client just needs to be initialized with a custom
// LibP2P host-based transport to perform requests to such server:
//
//	tr := &http.Transport{}
//	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))
//	client := &http.Client{Transport: tr}
//	res, err := client.Get("libp2p://Qmaoi4isbcTbFfohQyn28EiYM5CDWQx9QRCjDh3CTeiY7P/hello")
//	...
//
// In the example above, the client registers a "libp2p" protocol for which the
// custom transport is used. It can still perform regular "http" requests. The
// protocol name used is arbitraty and non standard.
//
// Note that LibP2P hosts cannot dial to themselves, so there is no possibility
// of using the same host as server and as client.
package p2phttp

import protocol "github.com/libp2p/go-libp2p-protocol"

// P2PProtocol is used to tag and identify streams
// handled by go-libp2p-http
var P2PProtocol protocol.ID = "/libp2p-http"

// Network is the name identifying the network to which
// go-libp2p-http addresses belong.
var Network = "p2p"
