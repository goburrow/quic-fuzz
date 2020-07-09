package quic

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

	for i := 0; i < 5; i++ {
		if client.IsEstablished() && server.IsEstablished() {
			return
		}
		n, err := client.Read(buf)
		if err != nil {
			panic(err)
		}
		n, err = server.Write(buf[:n])
		if err != nil {
			panic(err)
		}
		n, err = server.Read(buf)
		if err != nil {
			panic(err)
		}
		n, err = client.Write(buf[:n])
		if err != nil {
			panic(err)
		}
	}
	if !client.IsEstablished() || !server.IsEstablished() {
		panic("connection not established")
	}
	return
}
