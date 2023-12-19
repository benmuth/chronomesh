package network

import (
	"encoding/gob"
	"net"
)

const Localhost = "127.0.0.1"

type Sender interface {
	Send([]byte) error
}

type Client struct {
	// remoteAddress *net.TCPAddr
	// remotePort string
	id   uint64
	conn *net.TCPConn
}

func NewClient(id uint64, addr string, port string) (Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr+":"+port)
	if err != nil {
		return Client{}, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return Client{}, err
	}
	return Client{
		// remoteAddress: tcpAddr,
		id:   id,
		conn: conn,
	}, nil
}

func (c Client) Send(payload any) error {
	enc := gob.NewEncoder(c.conn)
	if err := enc.Encode(&payload); err != nil {
		return err
	}
	return nil
}
