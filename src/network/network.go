package network

import (
	"encoding/gob"
	"fmt"
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
	Conn *net.TCPConn
}

type Server struct {
	id       uint64
	Listener net.Listener
	Conn     *net.Conn
}

// Callers are responsible for closing the connection
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
		Conn: conn,
	}, nil
}

func (c Client) SendTimeEntry(entry timer.TimeEntry) error {
	enc := gob.NewEncoder(c.Conn)
	if err := enc.Encode(&entry); err != nil {
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

func (s Server) WaitForTimeEntry(ec chan timer.TimeEntry) {
	// cc := make(chan net.Conn)

	conn, err := s.Listener.Accept()
	if err != nil {
		ec <- timer.TimeEntry{ErrMessage: err.Error()}
	}
	fmt.Println("got connection")
	// cc <- conn
	defer conn.Close()

	dec := gob.NewDecoder(conn)

	var entry timer.TimeEntry
	if err := dec.Decode(&entry); err != nil {
		ec <- timer.TimeEntry{ErrMessage: err.Error()}
	}
	ec <- entry
	fmt.Println("decoded message")
}
