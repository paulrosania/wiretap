package wiretap

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

const (
	CRLF = "\r\n"
)

type TestClient struct {
	t                *testing.T
	Network, Address string
	Terminator       string
	r                *bufio.Reader
	conn             net.Conn
}

func NewTestClient(t *testing.T, network, address string) *TestClient {
	return &TestClient{
		t:          t,
		Network:    network,
		Address:    address,
		Terminator: CRLF,
	}
}

func (c *TestClient) connect() {
	if c.conn != nil {
		return
	}

	conn, err := net.Dial(c.Network, c.Address)
	if err != nil {
		c.t.Fatal(err)
	}
	c.conn = conn
	c.r = bufio.NewReader(conn)
}

func (c *TestClient) Expect(expect string) {
	c.connect()

	actual, err := c.r.ReadString('\n')
	if err != nil {
		c.t.Fatal(err)
	}
	l := len(actual) - len(c.Terminator)
	c.t.Log("S:", actual[:l])

	expect += c.Terminator

	if actual != expect {
		c.t.Fatalf("Expected %q, got %q", expect, actual)
	}

	return
}

func (c *TestClient) send(s string) {
	c.connect()

	c.t.Log("C:", s)
	s += c.Terminator

	c.conn.Write([]byte(s))
}

func (c *TestClient) Send(args ...interface{}) {
	s := fmt.Sprintln(args...)
	l := len(s) - 1 // Trim '\n'
	s = s[:l]
	c.send(s)
}

func (c *TestClient) Sendf(format string, args ...interface{}) { c.send(fmt.Sprintf(format, args...)) }

func (c *TestClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
