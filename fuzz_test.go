package quic

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/goburrow/quic/transport"
)

// Test and log handshake transactions
func TestHandshake(t *testing.T) {
	cid := make([]byte, transport.MaxCIDLength)
	client := newClient(cid)
	server := newServer(cid)
	err := logHandshake(client, server)
	if err != nil {
		t.Fatal(err)
	}
	err = logStream(client, server)
	if err != nil {
		t.Fatal(err)
	}
}

func logHandshake(client, server *transport.Conn) error {
	for i := 0; i < 10; i++ {
		if client.IsEstablished() && server.IsEstablished() {
			return nil
		}
		n, err := client.Read(buf)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fmt.Sprintf("client-%d", i), buf[:n], 0644)
		if err != nil {
			return err
		}
		n, err = server.Write(buf[:n])
		if err != nil {
			return err
		}
		n, err = server.Read(buf)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fmt.Sprintf("server-%d", i), buf[:n], 0644)
		n, err = client.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	if !client.IsEstablished() || !server.IsEstablished() {
		return fmt.Errorf("connection not established")
	}
	return nil
}

func logStream(client, server *transport.Conn) error {
	msg := []byte("hello")
	cs, err := client.Stream(2)
	if err != nil {
		return err
	}
	cs.Write(msg)
	n, err := client.Read(buf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("client-stream-0", buf[:n], 0644)
	if err != nil {
		return err
	}
	ss, err := server.Stream(3)
	if err != nil {
		return err
	}
	ss.Write(msg)
	n, err = server.Read(buf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("server-stream-0", buf[:n], 0644)
	if err != nil {
		return err
	}
	return nil
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
