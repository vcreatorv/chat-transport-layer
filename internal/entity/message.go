package entity

import "time"

type Message struct {
	Received int
	Total    int
	Last     time.Time
	Username string
	Segments []string
}
