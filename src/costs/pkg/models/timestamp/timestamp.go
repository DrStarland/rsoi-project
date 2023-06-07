package timestamp

import "time"

type Timestamp string

func Now() string {
	return time.Now().Format("2006-01-02T15:04:05.000")
}
