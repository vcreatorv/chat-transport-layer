package entity

import "time"

type Segment struct {
	SegmentNumber  int       `json:"segment_number"`
	TotalSegments  int       `json:"total_segments"`
	Username       string    `json:"username"`
	SendTime       time.Time `json:"send_time"`
	SegmentPayload string    `json:"payload"`
}
