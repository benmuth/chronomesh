package timer

import (
	"time"
)

type TimeEntry struct {
	Category   string
	Start, End time.Time
	Ended      bool
}

// func Start() {
// 	start := time.Now()
// 	print(start)
// }
