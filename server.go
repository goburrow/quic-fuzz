package quicfuzz

import (
	"crypto/tls"

	"github.com/goburrow/quic/transport"
)

var serverConfig = newServerConfig()

// FuzzServerInitial runs fuzzing server handling initial packet.
func FuzzServerInitial(b []byte) int {
	h := transport.Header{}
	_, err := h.Decode(b, transport.MaxCIDLength)
	if err != nil {
		return 0
	}
	conn, err := transport.Accept(h.DCID, nil, serverConfig)
	if err != nil {
		return 0
	}
	n, err := conn.Write(b)
	conn.Read(buf)
	if err != nil || n == 0 {
		return 0
	}
	return 1
}

// FuzzServer runs fuzzing connected server connection.
func FuzzServer(b []byte) int {
	peer, conn := newEndpoint()
	_, err1 := conn.Write(b)
	_, err2 := conn.Write(peer.BuildPacket(b))
	conn.Read(buf)
	if err1 != nil || err2 != nil {
		return 0
	}
	return 1
}

func newServer(scid []byte) *transport.Conn {
	conn, err := transport.Accept(scid, nil, serverConfig)
	if err != nil {
		panic(err)
	}
	return conn
}

func newServerConfig() *transport.Config {
	const certPEM = `-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`
	const keyPEM = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIrYSSNQFaA2Hwf1duRSxKtLYX5CB04fSeQ6tF1aY/PuoAoGCCqGSM49
AwEHoUQDQgAEPR3tU2Fta9ktY+6P9G0cWO+0kETA6SFs38GecTyudlHz6xvCdz8q
EKTcWGekdmdDPsHloRNtsiCa697B2O9IFA==
-----END EC PRIVATE KEY-----`
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		panic(err)
	}
	c := newConfig()
	c.TLS.Certificates = []tls.Certificate{cert}
	return c
}
