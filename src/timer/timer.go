package timer

import (
	"time"
)

type TimeEntry struct {
	Id         uint64
	Category   string
	Start, End time.Time
	Ended      bool
	ErrMessage string
}

type TimeEntries []TimeEntry

func (te TimeEntries) Merge(te2 TimeEntries) TimeEntries {
	return append(te, te2...)
}
