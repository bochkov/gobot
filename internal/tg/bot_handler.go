package tg

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) BotHandler(w http.ResponseWriter, req *http.Request) {
	var upd Update
	if err := json.NewDecoder(req.Body).Decode(&upd); err != nil {
		return
	}
	log.Printf("%+v", upd.Message)
	token := chi.URLParam(req, "token")
	if h.shouldAnswer(upd.Message) {
		go h.sendAnswer(token, upd.Message)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) shouldAnswer(msg *Message) bool {
	return msg != nil &&
		(strings.HasPrefix(msg.Text, "@resnyx") || strings.HasPrefix(msg.Text, "/"))
}

func (h *Handler) sendAnswer(token string, msg *Message) {
	methods := h.Service.GetAnswers(msg)
	for _, method := range methods {
		res, err := h.Service.Execute(method, token)
		if err != nil {
			log.Print(err)
			return
		}
		if res.Ok {
			log.Printf("response: %+v", res.Result)
		} else {
			log.Printf("response: %d %s", res.ErrorCode, res.Description)
		}
	}
}
