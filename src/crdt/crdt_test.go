package crdt

import (
	"testing"
	"time"

	"github.com/benmuth/time-tracker/src/timer"
	"github.com/google/go-cmp/cmp"
)

func TestLWWCommutativity(t *testing.T) {
	tests := []struct {
		name    string
		lwwSet1 LWWElementSet
		lwwSet2 LWWElementSet
	}{
		{
			name: "merge commutativity",
			lwwSet1: LWWElementSet{
				Bias: Add,
				AddSet: EntrySet{timer.TimeEntry{
					Id:         1,
					Category:   "code",
					Start:      time.Now(),
					End:        time.Now(),
					Ended:      true,
					ErrMessage: "",
				}: struct{}{}},
				RemoveSet: EntrySet{
					timer.TimeEntry{
						Id:         2,
						Category:   "",
						Start:      time.Now().Add(time.Second),
						End:        time.Now().Add(time.Second),
						Ended:      true,
						ErrMessage: "",
					}: struct{}{},
				},
			},
			lwwSet2: LWWElementSet{
				Bias:      Add,
				AddSet:    newEntrySet(1, 1),
				RemoveSet: newEntrySet(1, 2),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {
			mergeForward := tc.lwwSet1.Merge(tc.lwwSet2)
			mergeBackward := tc.lwwSet2.Merge(tc.lwwSet1)
			if diff := cmp.Diff(mergeForward, mergeBackward); diff != "" {
				t.Errorf("Non commutative merge. diff: %s", diff)
			}
		})
	}
}
