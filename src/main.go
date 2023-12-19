package main

import (
	"fmt"
	"log"
	"math/rand"
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
	server, err := network.NewServer(2, "", ServerPort)
	if err != nil {
		log.Fatalf("failed to listen for connections")
	}
	defer server.Listener.Close()

	var entry timer.TimeEntry
	entry, err = server.ReceiveTimeEntry(entry)
	if err != nil {
		log.Fatalf("failed to receive entry: %s", err)
	}
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
