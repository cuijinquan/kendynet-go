package util

import "time"

func SystemMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
