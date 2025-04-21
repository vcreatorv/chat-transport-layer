package entity

import "time"

type ReceiveRequest struct {
	Username string    `json:"username"`
	Data     string    `json:"data"`
	SendTime time.Time `json:"send_time"`
	Error    string    `json:"error"`
}
