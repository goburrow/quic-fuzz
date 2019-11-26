package quic

import (
	"crypto/tls"
	"time"

	"github.com/goburrow/quic/transport"
)

func stubTime() time.Time {
	return time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
}

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
		Time: stubTime,
	}
	return c
}
