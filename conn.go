package p2phttp

import (
	"net"
	"time"

	pnet "github.com/libp2p/go-libp2p-net"
)

// Conn is an implementation of net.Conn which wraps
// libp2p streams.
type Conn struct {
	s pnet.Stream
}

// NewConn creates a Conn given a libp2p stream
func NewConn(s pnet.Stream) net.Conn {
	return &Conn{s}
}

// Read reads data from the connection.
func (c *Conn) Read(b []byte) (n int, err error) {
	return c.s.Read(b)
}

// Write writes data to the connection.
func (c *Conn) Write(b []byte) (n int, err error) {
	return c.s.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() error {
	return c.s.Close()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return &Addr{c.s.Conn().LocalPeer()}
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return &Addr{c.s.Conn().RemotePeer()}
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
// See https://golang.org/pkg/net/#Conn for more details.
func (c *Conn) SetDeadline(t time.Time) error {
	return c.s.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls.
// A zero value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.s.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.s.SetWriteDeadline(t)
}
