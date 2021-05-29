package quicfuzz

import (
	"crypto/tls"
	"time"

	"github.com/goburrow/quic/transport"
)

var (
	cid = make([]byte, transport.MaxCIDLength)
	buf = make([]byte, 1400)
)

type stubRand struct{}

func (stubRand) Read(b []byte) (int, error) {
	for i := range b {
		b[i] = 0
	}
	return len(b), nil
}

func newConfig() *transport.Config {
	c := transport.NewConfig()
	c.TLS = &tls.Config{
		Rand: stubRand{},
		Time: func() time.Time {
			return time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
		},
	}
	return c
}

func newEndpoint() (client, server *transport.Conn) {
	client = newClient(cid)
	server = newServer(cid)

	i := 0
	for client.ConnectionState() != transport.StateActive ||
		server.ConnectionState() != transport.StateActive {
		n, err := client.Read(buf)
		if err != nil {
			panic(err)
		}
		if n > 0 {
			n, err = server.Write(buf[:n])
			if err != nil {
				panic(err)
			}
		}
		n, err = server.Read(buf)
		if err != nil {
			panic(err)
		}
		if n > 0 {
			n, err = client.Write(buf[:n])
			if err != nil {
				panic(err)
			}
		}
		i++
		if i > 5 {
			panic("connections not established")
		}
	}
	return
}
