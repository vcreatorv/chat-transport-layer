package utils

import (
	"TransportLayer/internal/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CodeURL         = "http://172.20.10.5:3050/api/code"
	MarsReceiveURL  = "http://172.20.10.2:8010/receive"
	EarthReceiveURL = "http://172.20.10.10:8001/receive"
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

	fmt.Println(resp.Status)
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
