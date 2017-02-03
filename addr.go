package p2phttp

import peer "github.com/libp2p/go-libp2p-peer"

// Addr implements net.Addr and holds a libp2p peer ID.
type Addr struct{ id peer.ID }

func (a *Addr) Network() string { return "libp2p" }
func (a *Addr) String() string  { return a.id.Pretty() }
