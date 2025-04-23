package http

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/delivery/kafka"
	"TransportLayer/internal/entity"
	"TransportLayer/internal/usecase"
	"TransportLayer/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type MessageHandler struct {
	msgUC    usecase.MessageService
	cfg      config.KafkaConfig
	producer kafka.Producer
}

func NewMessageHandler(msgUC usecase.MessageService, cfg config.KafkaConfig, producer kafka.Producer) *MessageHandler {
	return &MessageHandler{
		msgUC:    msgUC,
		cfg:      cfg,
		producer: producer,
	}
}

func (h *MessageHandler) Configure(r *mux.Router) {
	r.HandleFunc("/send", h.HandleSend).Methods("POST")
	r.HandleFunc("/transfer", h.HandleTransfer).Methods("POST")
}

func (h *MessageHandler) HandleSend(w http.ResponseWriter, r *http.Request) {
	var message entity.SendRequest
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, entity.ErrBadRequest, http.StatusBadRequest)
		return
	}

	segments, err := h.msgUC.SegmentMessage(message)
	if err != nil {
		http.Error(w, entity.ErrInternal, http.StatusInternalServerError)
		return
	}

	for _, segment := range segments {
		go utils.CodeSegment(segment)
		fmt.Printf("сегмент отправился: %+v\n", segment)
	}
}

func (h *MessageHandler) HandleTransfer(w http.ResponseWriter, r *http.Request) {
	var segment entity.Segment
	if err := json.NewDecoder(r.Body).Decode(&segment); err != nil {
		http.Error(w, entity.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := h.producer.WriteToKafka(segment); err != nil {
		http.Error(w, entity.ErrInternal, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
