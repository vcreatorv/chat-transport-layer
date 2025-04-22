package utils

import (
	"TransportLayer/internal/entity"
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	CODE_URL = "http://localhost:8000/api/code"
)

func SendSegment(body *entity.Segment) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", CODE_URL, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
