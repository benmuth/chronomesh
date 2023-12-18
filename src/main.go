package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/benmuth/time-tracker/src/timer"
)

const (
	ServerPort = "12345"
)

func main() {
	if os.Args[1] == "server" {
		server()
	} else if os.Args[1] == "client" {
		client()
	} else {
		fmt.Printf("Argument not recognized: %s\n", os.Args[1])
	}
}

func server() {
	listener, err := net.Listen("tcp", ":"+ServerPort)
	if err != nil {
		log.Fatalf("Error occurred: %s", err.Error())
	}
	defer listener.Close()

	fmt.Printf("Server is listening for connections on port %s...\n", ServerPort)

	fmt.Printf("waiting and listening!\n")
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("Error accepting connection: %s\n", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("connected!")

	// b := make([]byte, 100)

	var entry timer.TimeEntry
	// var timeStart time.Time
	dec := gob.NewDecoder(conn)
	if err := dec.Decode(&entry); err != nil {
		log.Fatalf("couldn't read entry from connection: %s", err)
	}
	fmt.Println(entry)

	// n, err := conn.Read(b)
	// if err != nil {
	// 	fmt.Printf("failed to read from connection\n")
	// }
	// fmt.Printf("%d bytes read\n", n)
	// fmt.Printf("message: %s\n", b)
}

func client() {
	// target_ip := os.Getenv("MY_IP")
	target_ip := "127.0.0.1"
	fmt.Println("IP: ", target_ip)
	addr, err := net.ResolveTCPAddr("tcp", target_ip+":"+ServerPort)
	if err != nil {
		log.Fatalf("couldn't resolve server IP: %s", addr)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalf("couldn't connect to server: %s", err)
	}
	defer conn.Close()

	entry := timer.TimeEntry{Category: "code", Start: time.Now(), End: time.Now(), Ended: true}

	// entryVal := reflect.ValueOf(&entry)
	// fmt.Printf("can set entry: %v\n", entryVal.Elem().CanSet())

	enc := gob.NewEncoder(conn)
	if err := enc.Encode(&entry); err != nil {
		log.Fatalf("couldn't encode entry: %s", err)
	}
	// n, err := conn.Write([]byte("hello!"))
	// if err != nil {
	// 	log.Fatalf("couldn't write to connection")
	// }
	// fmt.Printf("%v bytes written\n", n)
}