package utils

import (
	"TransportLayer/internal/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CodeURL         = "http://192.168.204.34:3050/api/code"
	MarsReceiveURL  = "http://192.168.204.224:8010/receive"
	EarthReceiveURL = "http://192.168.204.153:8001/receive"
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
	var ReceiveURL = MarsReceiveURL
	if body.Error != "" {
		ReceiveURL = EarthReceiveURL
	}

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
