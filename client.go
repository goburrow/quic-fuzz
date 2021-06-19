package quicfuzz

import (
	"github.com/goburrow/quic/transport"
)

var clientConfig = newClientConfig()

// FuzzClientInitial runs fuzzing client handling initial packet.
func FuzzClientInitial(b []byte) int {
	conn := newClient(cid)
	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	n, err := conn.Write(b)
	if err != nil || n == 0 {
		return 0
	}
	conn.Read(buf)
	return 1
}

// FuzzClient runs fuzzing connected client connection.
func FuzzClient(b []byte) int {
	conn, peer := newEndpoint()
	_, err1 := conn.Write(b)
	_, err2 := conn.Write(peer.BuildPacket(b))
	conn.Read(buf)
	if err1 != nil || err2 != nil {
		return 0
	}
	return 1
}

func newClient(scid []byte) *transport.Conn {
	conn, err := transport.Connect(scid, scid, clientConfig)
	if err != nil {
		panic(err)
	}
	return conn
}

func newClientConfig() *transport.Config {
	c := newConfig()
	c.TLS.InsecureSkipVerify = true
	return c
}
