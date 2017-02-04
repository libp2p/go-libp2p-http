package p2phttp

import peer "github.com/libp2p/go-libp2p-peer"

// Addr implements net.Addr and holds a libp2p peer ID.
type Addr struct{ id peer.ID }

// Network returns the name of the network that this address belongs to
// (libp2p).
func (a *Addr) Network() string { return "libp2p" }

// String returns the peer ID of this address in string form
// (B58-encoded).
func (a *Addr) String() string { return a.id.Pretty() }
