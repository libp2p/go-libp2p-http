package p2phttp

import (
	"context"
	"net"
	"time"

	host "github.com/libp2p/go-libp2p-host"
	pnet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
)

// Listener is an implementation of net.Listener which handles
// http-tagged streams from a libp2p connection.
// A listener can be built with Listen()
type Listener struct {
	host     host.Host
	ctx      context.Context
	cancel   func()
	streamCh chan pnet.Stream
}

// Accept returns a connection from this listener. It blocks if there
// are no connections.
func (l *Listener) Accept() (net.Conn, error) {
	select {
	case s := <-l.streamCh:
		return NewConn(s), nil
	case <-l.ctx.Done():
		return nil, l.ctx.Err()
	}
}

// Close terminates this listener. It will no longer handle any
// incoming streams
func (l *Listener) Close() error {
	l.cancel()
	l.host.RemoveStreamHandler(P2PProtocol)
	return nil
}

// Addr returns the address for this listener, which is its libp2p Peer ID.
func (l *Listener) Addr() net.Addr {
	return &Addr{l.host.ID()}
}

// Listen creates a new listener ready to accept http-tagged streams
// received by a host.
func Listen(h host.Host) (net.Listener, error) {
	ctx, cancel := context.WithCancel(context.Background())

	l := &Listener{
		host:     h,
		ctx:      ctx,
		cancel:   cancel,
		streamCh: make(chan pnet.Stream),
	}

	h.SetStreamHandler(P2PProtocol, func(s pnet.Stream) {
		select {
		case l.streamCh <- s:
		case <-ctx.Done():
			s.Close()
		}
	})

	return l, nil
}

// Dial opens a stream to the destination address
// (which should parseable to a peer ID) using the given
// host and returns it.
func Dial(h host.Host, address string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	pid, err := peer.IDB58Decode(address)
	if err != nil {
		return nil, err
	}

	s, err := h.NewStream(ctx, pid, P2PProtocol)
	if err != nil {
		return nil, err
	}
	return NewConn(s), nil
}
