package entity

import "time"

type SendRequest struct {
	Username string    `json:"username"`
	Data     string    `json:"data"`
	SendTime time.Time `json:"send_time"`
}
