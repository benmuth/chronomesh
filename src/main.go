package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/benmuth/time-tracker/src/network"
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
	client, err := network.NewClient(1, network.Localhost, ServerPort)
	if err != nil {
		log.Fatalf("failed to initialize client %s", err)
	}

	entry := timer.TimeEntry{Category: "code", Start: time.Now(), End: time.Now(), Ended: true, Id: rand.Uint64()}

	if err := client.Send(entry); err != nil {
		log.Fatalf("couldn't send entry: %s", err)
	}
}
