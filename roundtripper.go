package p2phttp

import (
	"bufio"
	"net/http"

	host "github.com/libp2p/go-libp2p-host"
)

// RoundTripper implemenets http.RoundTrip and can be used as
// custom transport with Go http.Client.
type RoundTripper struct {
	h host.Host
}

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
func (rt *RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	addr := r.Host
	if addr == "" {
		addr = r.URL.Host
	}

	conn, err := Dial(rt.h, r.URL.Host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = r.Write(conn)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	return http.ReadResponse(reader, r)
}

// NewTransport returns a new RoundTripper which uses the provided
// libP2P host to perform an http request and obtain the response.
//
// The typical use case for NewTransport is to register the "libp2p"
// protocol with a Transport, as in:
//     t := &http.Transport{}
//     t.RegisterProtocol(p2phttp.P2PProtocol, p2phttp.NewTransport(host))
//     c := &http.Client{Transport: t}
//     res, err := c.Get("libp2p://Qmaoi4isbcTbFfohQyn28EiYM5CDWQx9QRCjDh3CTeiY7P/index.html")
//     ...
func NewTransport(h host.Host) *RoundTripper {
	return &RoundTripper{h}
}
