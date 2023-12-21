package network

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/benmuth/time-tracker/src/timer"
	"github.com/google/go-cmp/cmp"
)

const ServerPort = "12345"

func TestSendAndReceive(t *testing.T) {
	tests := []struct {
		name  string
		entry timer.TimeEntry
	}{
		{
			name: "single exchange",
			entry: timer.TimeEntry{
				Category: "code",
				Start:    time.Now(),
				End:      time.Now(),
				Ended:    true,
				Id:       rand.Uint64(),
			},
		},
	}

	server, err := NewServer(2, "", ServerPort)
	if err != nil {
		log.Fatalf("failed to listen for connections")
	}
	defer server.Listener.Close()

	ec := make(chan timer.TimeEntry)
	go server.WaitForTimeEntry(ec)

	client, err := NewClient(1, Localhost, ServerPort)
	if err != nil {
		log.Fatalf("failed to initialize client %s", err)
	}

	for _, tc := range tests {
		if err := client.SendTimeEntry(tc.entry); err != nil {
			log.Fatalf("couldn't send entry: %s", err)
		}

		receivedEntry := <-ec
		if receivedEntry.ErrMessage != "" {
			t.Fatalf("Failed to receive entry: %s", receivedEntry.ErrMessage)
		}

		if diff := cmp.Diff(tc.entry, receivedEntry); diff != "" {
			t.Errorf("sent and received structs don't match. mismatch (-want +got):\n %s", diff)
		}
	}
}
