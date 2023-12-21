package timer

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestMergeSymmetry(t *testing.T) {
	tests := []struct {
		name          string
		localEntries  TimeEntries
		remoteEntries TimeEntries
		wantLength    int
	}{
		{
			name:          "Bidirectional symmetry",
			localEntries:  generateEntries(2),
			remoteEntries: generateEntries(2),
		},
	}

	for _, tc := range tests {
		localToRemote := tc.localEntries.Merge(tc.remoteEntries)
		remoteToLocal := tc.remoteEntries.Merge(tc.localEntries)
		if diff := cmp.Diff(localToRemote, remoteToLocal); diff != "" {
			t.Errorf("Asymmetric merge: %s", diff)
		}
	}
}

func generateEntries(n int) TimeEntries {
	entries := make([]TimeEntry, 0)
	for i := 0; i < n; i++ {
		newEntry := TimeEntry{
			Id:         rand.Uint64(),
			Category:   "code",
			Start:      time.Now(),
			End:        time.Now(),
			Ended:      true,
			ErrMessage: "",
		}
		entries = append(entries, newEntry)
	}
	return entries
}
