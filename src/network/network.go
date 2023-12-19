package network

import (
	"encoding/gob"
	"net"

	"github.com/benmuth/time-tracker/src/timer"
)

const Localhost = "127.0.0.1"

// type Sender interface {
// 	Send([]byte) error
// }

// type Receiver interface {
// 	ReceiveTimeEntry(timer.TimeEntry) (timer.TimeEntry, error)
// }

type Client struct {
	id   uint64
	conn *net.TCPConn
}

type Server struct {
	id       uint64
	Listener net.Listener
	conn     *net.Conn
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

// It's the caller's responsibility to close the listener.
// Call Server.Listener.Accept to accept an incoming connection.
func NewServer(id uint64, addr string, port string) (Server, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return Server{}, err
	}

	return Server{
		id: id, Listener: listener,
	}, nil
}

func (s Server) ReceiveTimeEntry(entry timer.TimeEntry) (timer.TimeEntry, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		return timer.TimeEntry{}, err
	}
	defer conn.Close()

	dec := gob.NewDecoder(conn)
	if err := dec.Decode(&entry); err != nil {
		return timer.TimeEntry{}, err
	}
	return entry, nil
}
