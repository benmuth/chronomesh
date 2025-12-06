package crdt

import (
	"time"

	"github.com/benmuth/time-tracker/src/timer"
)

const (
	Add = iota
	Remove
)

//	type TimeEntry struct {
//		Id         uint64
//		Category   string
//		Start, End time.Time
//		Ended      bool
//		ErrMessage string
//	}
type LWWElementSet struct {
	Bias      int
	AddSet    EntrySet
	RemoveSet EntrySet
}

type EntrySet map[timer.TimeEntry]struct{}

// newEntrySet creates a entry set with sequential timestamps
func newEntrySet(n int, timeOffset int) EntrySet {
	es := make(EntrySet)
	for i := 0; i < n; i++ {
		te := timer.NewTimeEntry(time.Now().Add(time.Duration(i)*time.Second), time.Now().Add(time.Duration(i)*time.Second+1))
		es[te] = struct{}{}
	}
	return es
}

// TODO: make this merge idempotent, associative, and commutative
func (lww *LWWElementSet) Merge(lwwRemote LWWElementSet) LWWElementSet {
	newLWWSet := LWWElementSet{Bias: Add, AddSet: make(EntrySet), RemoveSet: make(EntrySet)}
	for k := range lww.AddSet {
		newLWWSet.AddSet[k] = struct{}{}
	}
	for k := range lww.RemoveSet {
		newLWWSet.RemoveSet[k] = struct{}{}
	}
	for k := range lwwRemote.AddSet {
		newLWWSet.AddSet[k] = struct{}{}
	}
	for k := range lwwRemote.RemoveSet {
		newLWWSet.RemoveSet[k] = struct{}{}
	}
	return *lww
}
