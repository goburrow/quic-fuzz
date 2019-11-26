package quic

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/goburrow/quic/transport"
)

func logHandshake(client, server *transport.Conn) error {
	b := make([]byte, 1200)
	for i := 0; i < 10; i++ {
		if client.IsEstablished() && server.IsEstablished() {
			break
		}
		n, err := client.Read(b)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fmt.Sprintf("client-%d", i), b[:n], 0644)
		if err != nil {
			return err
		}
		n, err = server.Write(b[:n])
		if err != nil {
			return err
		}
		n, err = server.Read(b)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fmt.Sprintf("server-%d", i), b[:n], 0644)
		n, err = client.Write(b[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

// Test and log handshake transactions
func TestHandshake(t *testing.T) {
	cid := make([]byte, transport.MaxCIDLength)
	client, err := transport.Connect(cid, clientConfig)
	if err != nil {
		t.Fatal(err)
	}
	server, err := transport.Accept(cid, nil, serverConfig)
	if err != nil {
		t.Fatal(err)
	}
	err = logHandshake(client, server)
	if err != nil {
		t.Fatal(err)
	}
	if !client.IsEstablished() || !server.IsEstablished() {
		t.Fatalf("connection not established")
	}
}

func TestServerInitial(t *testing.T) {
	b, err := ioutil.ReadFile("client-0")
	if err != nil {
		t.Fatal(err)
	}
	n := FuzzServerInitial(b)
	if n != 1 {
		t.Fatalf("fuzz: %d", n)
	}
}
