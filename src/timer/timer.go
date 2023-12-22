package timer

import (
	"math/rand"
	"time"
)

type TimeEntry struct {
	Id         uint64
	Category   string
	Start, End time.Time
	Ended      bool
	ErrMessage string
}

func NewTimeEntry(start, end time.Time) TimeEntry {
	return TimeEntry{
		Id:         rand.Uint64(),
		Category:   "",
		Start:      start,
		End:        end,
		Ended:      false,
		ErrMessage: "",
	}
}

type TimeEntries []TimeEntry

func (te TimeEntries) Merge(te2 TimeEntries) TimeEntries {
	return append(te, te2...)
}
