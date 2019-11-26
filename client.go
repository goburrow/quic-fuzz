package quic

import (
	"github.com/goburrow/quic/transport"
)

var clientConfig = newClientConfig()
var buf = make([]byte, 1400)

func FuzzClientInitial(b []byte) int {
	conn, err := transport.Connect([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, clientConfig)
	if err != nil {
		panic(err)
	}
	_, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}
	n, err := conn.Write(b)
	if err != nil || n == 0 {
		return 0
	}
	return 1
}

func newClientConfig() *transport.Config {
	c := newConfig()
	c.TLS.InsecureSkipVerify = true
	return c
}
