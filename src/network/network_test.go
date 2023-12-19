package network

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/benmuth/time-tracker/src/timer"
)

const ServerPort = "12345"

func TestSendAndReceive(t *testing.T) {
	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		id    uint64
	}{
		{
			name:  "single exchange",
			start: time.Now(),
			end:   time.Now(),
			id:    rand.Uint64(),
		},
	}

	server, err := NewServer(2, "", ServerPort)
	if err != nil {
		log.Fatalf("failed to listen for connections")
	}
	defer server.Listener.Close()

	fmt.Println("1")
	// var entry timer.TimeEntry
	ec := make(chan timer.TimeEntry)
	go server.WaitForTimeEntry(ec)

	fmt.Println("2")
	client, err := NewClient(1, Localhost, ServerPort)
	if err != nil {
		log.Fatalf("failed to initialize client %s", err)
	}

	for _, tc := range tests {
		entry := timer.TimeEntry{Category: "code", Start: tc.start, End: tc.end, Ended: true, Id: tc.id}

		if err := client.SendTimeEntry(entry); err != nil {
			log.Fatalf("couldn't send entry: %s", err)
		}

		receivedEntry := <-ec
		if receivedEntry.ErrMessage != "" {
			t.Fatalf("Failed to receive entry: %s", receivedEntry.ErrMessage)
		}
		if receivedEntry.Id != tc.id {
			t.Errorf("Ids don't match: want %d got %d\n", receivedEntry.Id, tc.id)
		}

		if receivedEntry.Start.Compare(tc.start) != 0 {
			diff := receivedEntry.Start.Sub(tc.start)
			t.Errorf("Start times don't match: want %v got %v\tdiff %v\n", receivedEntry.Start, tc.start, diff)

		}
		if receivedEntry.End.Compare(tc.end) != 0 {
			diff := receivedEntry.End.Sub(tc.end)
			t.Errorf("End times don't match: want %v got %v\tdiff %v\n", receivedEntry.End, tc.end, diff)
		}
	}
}
