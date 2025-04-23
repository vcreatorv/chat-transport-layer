package utils

import (
	"TransportLayer/internal/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CodeURL    = "http://localhost:8000/api/code"
	ReceiveURL = "http://localhost:3000/api/receive"
)

func CodeSegment(body *entity.Segment) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("ошибка сериализации сегмента")
		return
	}

	req, _ := http.NewRequest(http.MethodPost, CodeURL, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func SendMessage(body entity.ReceiveRequest) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, ReceiveURL, bytes.NewBuffer(reqBody))

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
